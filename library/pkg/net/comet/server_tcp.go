package comet

import (
    "context"
    "io"
    "net"
    "strings"
    "time"

    "github.com/znyh/middle-end/library/pkg/net/comet/internal/bucket"
    "github.com/znyh/middle-end/library/pkg/net/comet/internal/bufio"
    "github.com/znyh/middle-end/library/pkg/net/comet/internal/bytes"
    "github.com/znyh/middle-end/library/pkg/net/comet/internal/channel"
    "github.com/znyh/middle-end/library/pkg/net/comet/internal/proxy"
    xtime "github.com/znyh/middle-end/library/pkg/net/comet/internal/time"
    "github.com/znyh/middle-end/library/pkg/net/comet/proto"

    "github.com/go-kratos/kratos/pkg/log"
    "github.com/go-kratos/kratos/pkg/net/metadata"
    gproto "github.com/golang/protobuf/proto"
    "github.com/golang/protobuf/ptypes/empty"
    "github.com/google/uuid"
    grpcmd "google.golang.org/grpc/metadata"
)

const (
    maxInt = 1<<31 - 1
)

// StartTCP listen all tcp.bind and start accept connections.
func (s *Server) StartTCP(accept int) (err error) {
    var (
        bind     string
        listener *net.TCPListener
        addr     *net.TCPAddr
    )

    for _, bind = range s.c.TCP.Bind {
        if addr, err = net.ResolveTCPAddr("tcp", bind); err != nil {
            log.Error("net.ResolveTCPAddr(tcp, %s) error(%v)", bind, err)
            return
        }
        if listener, err = net.ListenTCP("tcp", addr); err != nil {
            log.Error("net.ListenTCP(tcp, %s) error(%v)", bind, err)
            return
        }
        log.Info("start tcp listen: %s", bind)
        // split N core accept
        for i := 0; i < accept; i++ {
            go s.acceptTCP(listener)
        }
    }
    return
}

func (s *Server) acceptTCP(lis *net.TCPListener) {
    var (
        conn *net.TCPConn
        err  error
        r    int
    )
    for {
        if conn, err = lis.AcceptTCP(); err != nil {
            // if listener close then return
            log.Error("listener.Accept(\"%s\") error(%v)", lis.Addr().String(), err)
            return
        }
        if err = conn.SetKeepAlive(s.c.Comet.KeepAlive); err != nil {
            log.Error("conn.SetKeepAlive() error(%v)", err)
            return
        }
        if err = conn.SetReadBuffer(s.c.Comet.Rcvbuf); err != nil {
            log.Error("conn.SetReadBuffer() error(%v)", err)
            return
        }
        if err = conn.SetWriteBuffer(s.c.Comet.Sndbuf); err != nil {
            log.Error("conn.SetWriteBuffer() error(%v)", err)
            return
        }
        go s.serveTCP(conn, r)
        if r++; r == maxInt {
            r = 0
        }
    }
}

func (s *Server) serveTCP(conn *net.TCPConn, r int) {
    defer func() {
        if err := s.recoveryServer(); err != nil {
            log.Error("serceTcp %v", err)
            conn.Close()
        }
    }()
    var (
        tr    = s.round.Timer(r)
        rp    = s.round.Reader(r)
        wp    = s.round.Writer(r)
        lAddr = conn.LocalAddr().String()
        rAddr = conn.RemoteAddr().String()
    )
    log.Info("start tcp serve \"%s\" with \"%s\"", lAddr, rAddr)
    var (
        err error
        hb  time.Duration
        p   *proto.Payload
        b   *bucket.Bucket
        trd *xtime.TimerData
        rb  = rp.Get()
        wb  = wp.Get()
        ch  = channel.NewChannel(s.c.Protocol.CliProto, s.c.Protocol.SvrProto)
        rr  = &ch.Reader
        wr  = &ch.Writer
    )
    ch.Reader.ResetBuffer(conn, rb.Bytes())
    ch.Writer.ResetBuffer(conn, wb.Bytes())
    //handshake
    uid := uuid.New().String()
    step := 0
    trd = tr.Add(time.Duration(s.c.Protocol.HandshakeTimeout), func() {
        s.disconnectChan <- uid
        conn.Close()
        log.Error("key: %s remoteIP: %s step: %d tcp handshake timeout", ch.Key, conn.RemoteAddr().String(), step)
    })
    //proxy
    if s.c.Protocol.Proxy {
        proxyHeader, err := proxy.Read(rr)
        //if err == proxy.ErrNoProxyProtocol {
        //	err = nil
        //}
        if err != nil {
            log.Error("proxy err %v", err)
        } else {
            lAddr = proxyHeader.LocalAddr().String()
            rAddr = proxyHeader.RemoteAddr().String()
            log.Info("proxy tcp serve \"%s\" with \"%s\"", lAddr, rAddr)
        }
    }
    ch.IP, _, _ = net.SplitHostPort(rAddr)
    step = 1
    //auth
    if s.c.Auth.Open {
        if p, err = ch.CliProto.Set(); err == nil {
            err = s.authTCP(conn, rr, p)
        }
    }
    ch.Key = uid
    hb = time.Duration(s.c.Protocol.HandshakeTimeout)
    b = s.GetBucket(ch.Key)
    b.Put(ch)
    md := metadata.MD{
        metadata.RemoteIP: ch.IP,
        metadata.Mid:      ch.Key,
    }
    newCtx := metadata.NewContext(context.Background(), md)
    ctx, cancel := context.WithCancel(newCtx)
    defer cancel()
    step = 2
    if err != nil {
        conn.Close()
        rp.Put(rb)
        wp.Put(wb)
        tr.Del(trd)
        log.Error("key: %s handshake failed error(%v)", ch.Key, err)
        return
    }
    trd.Key = ch.Key
    tr.Set(trd, hb)
    step = 3
    // hanshake ok start dispatch goroutine
    go s.dispatchTCP(conn, wr, wp, wb, ch)
    for {
        if p, err = ch.CliProto.Set(); err != nil {
            break
        }
        if err = p.ReadTCP(rr); err != nil {
            break
        }
        tr.Set(trd, time.Duration(s.c.Protocol.HeartbeatTimeout))
        if p.Type == int32(proto.Ping) {
            p.Type = int32(proto.Ping)
            p.Body = nil
            _metricServerReqCodeTotal.Inc("/Ping", "no_user", "0")
            step++
        } else {
            if err = s.Operate(ctx, p); err != nil {
                break
            }
        }
        //response为空的时候dispatchTCP不处理
        if p.Body != nil || p.Type == int32(proto.Ping) {
            ch.CliProto.SetAdv()
            ch.Signal()
        }
    }
    if err != nil && err != io.EOF && !strings.Contains(err.Error(), "closed") {
        log.Error("key: %s server tcp failed error(%v)", ch.Key, err)
    }
    s.disconnectChan <- uid
    b.Del(ch)
    tr.Del(trd)
    rp.Put(rb)
    conn.Close()
    ch.Close()
}

func (s *Server) dispatchTCP(conn *net.TCPConn, wr *bufio.Writer, wp *bytes.Pool, wb *bytes.Buffer, ch *channel.Channel) {
    defer func() {
        if err := s.recoveryServer(); err != nil {
            log.Error("dispatchTCP %v", err)
            conn.Close()
        }
    }()
    var (
        err    error
        finish bool
    )
    for {
        var p = ch.Ready()
        switch p {
        case proto.ProtoFinish:
            finish = true
            goto failed
        case proto.ProtoReady:
            // fetch message from svrbox(client send)
            for {
                if p, err = ch.CliProto.Get(); err != nil {
                    break
                }
                if p.Type == int32(proto.Ping) {
                    p.Type = int32(proto.Pong)
                    if err = p.WriteTCPHeart(wr); err != nil {
                        goto failed
                    }
                } else if p.Type == int32(proto.Response) {
                    if err = p.WriteTCP(wr); err != nil {
                        goto failed
                    }
                }
                p.Body = nil // avoid memory leak
                ch.CliProto.GetAdv()
            }
        default:
            // server send
            if err = p.WriteTCP(wr); err != nil {
                goto failed
            }
        }
        // only hungry flush response
        if err = wr.Flush(); err != nil {
            log.Error("Flush error(%v)", err)
            break
        }
    }
failed:
    if err != nil {
        log.Error("key: %s dispatch tcp error(%v)", ch.Key, err)
    }
    //s.disconnectChan <- ch.Key
    conn.Close()
    wp.Put(wb)
    // must ensure all channel message discard, for reader won't blocking Signal
    for !finish {
        finish = (ch.Ready() == proto.ProtoFinish)
    }
    return
}

func (s *Server) authTCP(conn *net.TCPConn, rr *bufio.Reader, p *proto.Payload) (err error) {
    reqBody := &proto.Body{}
    for {
        if err = p.ReadTCP(rr); err != nil {
            return
        }
        if p.Type == int32(proto.Request) {
            if err = gproto.Unmarshal(p.Body, reqBody); err != nil {
                return
            }
            if reqBody.Ops == proto.AuthOps {
                break
            } else {
                log.Error("tcp request ops(%d) not auth", reqBody.Ops)
            }
        }
    }
    token := string(reqBody.Data[:])
    ctx := grpcmd.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+token)
    _, err = s.atuhClient.VerifyToken(ctx, &empty.Empty{})
    return
}
