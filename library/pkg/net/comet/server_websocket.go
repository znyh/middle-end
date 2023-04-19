package comet

import (
	"context"
	"io"
	"net"
	"strings"
	"time"

	"github.com/go-kratos/kratos/pkg/log"
	"github.com/go-kratos/kratos/pkg/net/metadata"
	"github.com/google/uuid"
	"github.com/znyh/middle-end/library/pkg/net/comet/internal/bucket"
	"github.com/znyh/middle-end/library/pkg/net/comet/internal/bytes"
	"github.com/znyh/middle-end/library/pkg/net/comet/internal/channel"
	xtime "github.com/znyh/middle-end/library/pkg/net/comet/internal/time"
	"github.com/znyh/middle-end/library/pkg/net/comet/internal/websocket"
	"github.com/znyh/middle-end/library/pkg/net/comet/proto"
)

// InitWebsocket listen all tcp.bind and start accept connections.
func (s *Server) StartWebsocket(accept int) (err error) {
    var (
        bind     string
        listener *net.TCPListener
        addr     *net.TCPAddr
    )
    for _, bind = range s.c.Websocket.Bind {
        if addr, err = net.ResolveTCPAddr("tcp", bind); err != nil {
            log.Error("net.ResolveTCPAddr(tcp, %s) error(%v)", bind, err)
            return
        }
        if listener, err = net.ListenTCP("tcp", addr); err != nil {
            log.Error("net.ListenTCP(tcp, %s) error(%v)", bind, err)
            return
        }
        log.Info("start ws listen: %s", bind)
        // split N core accept
        for i := 0; i < accept; i++ {
            go s.acceptWebsocket(listener)
        }
    }
    return
}

// Accept accepts connections on the listener and serves requests
// for each incoming connection.  Accept blocks; the caller typically
// invokes it in a go statement.
func (s *Server) acceptWebsocket(lis *net.TCPListener) {
    var (
        conn *net.TCPConn
        err  error
        r    int
    )
    for {
        if conn, err = lis.AcceptTCP(); err != nil {
            // if listener close then return
            log.Error("listener.Accept(%s) error(%v)", lis.Addr().String(), err)
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
        go s.serveWebsocket(conn, r)
        if r++; r == maxInt {
            r = 0
        }
    }
}

func (s *Server) serveWebsocket(conn net.Conn, r int) {
    defer func() {
        if err := s.recoveryServer(); err != nil {
            log.Error("serveWebsocket %v", err)
            conn.Close()
        }
    }()
    var (
        // timer
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
        ch  = channel.NewChannel(s.c.Protocol.CliProto, s.c.Protocol.SvrProto)
        rr  = &ch.Reader
        wr  = &ch.Writer
        ws  *websocket.Conn // websocket
        req *websocket.Request
    )
    // reader
    ch.Reader.ResetBuffer(conn, rb.Bytes())
    // handshake
    uid := uuid.New().String()
    step := 0
    trd = tr.Add(time.Duration(s.c.Protocol.HandshakeTimeout), func() {
        s.disconnectChan <- uid
        // NOTE: fix close block for tls
        _ = conn.SetDeadline(time.Now().Add(time.Millisecond * 100))
        _ = conn.Close()
        log.Error("key: %s remoteIP: %s step: %d ws handshake timeout", ch.Key, conn.RemoteAddr().String(), step)
    })
    // websocket
    ch.IP, _, _ = net.SplitHostPort(conn.RemoteAddr().String())
    step = 1
    if req, err = websocket.ReadRequest(rr); err != nil || req.RequestURI != "/" {
        conn.Close()
        tr.Del(trd)
        rp.Put(rb)
        if err != io.EOF {
            log.Error("http.ReadRequest(rr) error(%v)", err)
        }
        return
    }
    ip := req.Header.Get("X-Forwarded-For")
    ip = strings.TrimSpace(strings.Split(ip, ",")[0])
    if ip == "" {
        ip = strings.TrimSpace(req.Header.Get("X-Real-Ip"))
    }
    if ip != "" {
        ch.IP = ip
    }
    // writer
    wb := wp.Get()
    ch.Writer.ResetBuffer(conn, wb.Bytes())
    step = 2
    if ws, err = websocket.Upgrade(conn, rr, wr, req); err != nil {
        conn.Close()
        tr.Del(trd)
        rp.Put(rb)
        wp.Put(wb)
        if err != io.EOF {
            log.Error("websocket.NewServerConn error(%v)", err)
        }
        return
    }
    // must not setadv, only used in auth
    step = 3
    ch.Mid = 0
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
    step = 4
    if err != nil {
        ws.Close()
        rp.Put(rb)
        wp.Put(wb)
        tr.Del(trd)
        if err != io.EOF && err != websocket.ErrMessageClose {
            log.Error("key: %s remoteIP: %s step: %d ws handshake failed error(%v)", ch.Key, conn.RemoteAddr().String(), step, err)
        }
        return
    }
    trd.Key = ch.Key
    tr.Set(trd, hb)
    // hanshake ok start dispatch goroutine
    step = 5
    go s.dispatchWebsocket(ws, wp, wb, ch)
    for {
        if p, err = ch.CliProto.Set(); err != nil {
            break
        }
        if err = p.ReadWebsocket(ws); err != nil {
            break
        }
        tr.Set(trd, time.Duration(s.c.Protocol.HeartbeatTimeout))
        if p.Op == int32(proto.Ping) {
            p.Op = int32(proto.Ping)
            p.Body = nil
            _metricServerReqCodeTotal.Inc("/Ping", "no_user", "0")
            step++
        } else {
            if err = s.Operate(ctx, p); err != nil {
                break
            }
        }
        ch.CliProto.SetAdv()
        ch.Signal()
    }
    if err != nil && err != io.EOF && err != websocket.ErrMessageClose && !strings.Contains(err.Error(), "closed") {
        log.Error("key: %s server ws failed error(%v)", ch.Key, err)
    }
    s.disconnectChan <- uid
    b.Del(ch)
    tr.Del(trd)
    ws.Close()
    ch.Close()
    rp.Put(rb)
}

// dispatch accepts connections on the listener and serves requests
// for each incoming connection.  dispatch blocks; the caller typically
// invokes it in a go statement.
func (s *Server) dispatchWebsocket(ws *websocket.Conn, wp *bytes.Pool, wb *bytes.Buffer, ch *channel.Channel) {
    defer func() {
        if err := s.recoveryServer(); err != nil {
            log.Error("dispatchWebsocket %v", err)
            ws.Close()
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
                    if err = p.WriteWebsocketHeart(ws); err != nil {
                        goto failed
                    }
                } else if p.Type == int32(proto.Response) {
                    if err = p.WriteWebsocket(ws); err != nil {
                        goto failed
                    }
                }
                p.Body = nil // avoid memory leak
                ch.CliProto.GetAdv()
            }
        default:
            if err = p.WriteWebsocket(ws); err != nil {
                goto failed
            }
        }
        // only hungry flush response
        if err = ws.Flush(); err != nil {
            break
        }
    }
failed:
    if err != nil && err != io.EOF && err != websocket.ErrMessageClose {
        log.Error("key: %s dispatch ws error(%v)", ch.Key, err)
    }
    ws.Close()
    wp.Put(wb)
    // must ensure all channel message discard, for reader won't blocking Signal
    for !finish {
        finish = (ch.Ready() == proto.ProtoFinish)
    }
}
