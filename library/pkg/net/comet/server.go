package comet

import (
    "context"
    "fmt"
    "log"
    "math"
    "os"
    "reflect"
    "runtime"
    "time"

    "github.com/znyh/middle-end/library/pkg/net/comet/internal/bucket"
    "github.com/znyh/middle-end/library/pkg/net/comet/internal/round"
    xtime "github.com/znyh/middle-end/library/pkg/net/comet/internal/time"
    "github.com/znyh/middle-end/library/pkg/net/comet/proto"

    "github.com/go-kratos/kratos/pkg/ecode"
    "github.com/go-kratos/kratos/pkg/net/rpc/warden"
    gproto "github.com/golang/protobuf/proto"
    "github.com/zhenjl/cityhash"
)

var (
    _defaultSerConf = &ServerConfig{
        Comet: &Comet{
            Sndbuf:       4096,
            Rcvbuf:       4096,
            KeepAlive:    false,
            Reader:       32,
            ReadBuf:      1024,
            ReadBufSize:  8192,
            Writer:       32,
            WriteBuf:     1024,
            WriteBufSize: 8192,
        },
        TCP: &TCP{
            Bind: []string{":3101"},
        },
        Websocket: &Websocket{
            Bind: []string{":3102"},
        },
        Protocol: &Protocol{
            Proxy:            false,
            Timer:            32,
            TimerSize:        2048,
            CliProto:         5,
            SvrProto:         10,
            HandshakeTimeout: xtime.Duration(time.Second * 15),
            HeartbeatTimeout: xtime.Duration(time.Second * 6),
        },
        Auth: &Auth{
            Open:  false,
            AppID: "auth",
        },
        Bucket: &Bucket{
            Size:    32,
            Channel: 1024,
        },
        ChanSize: &ChanSize{
            Push:       2048,
            Close:      1024,
            Disconnect: 1024,
        },
    }
    _abortIndex int8 = math.MaxInt8 / 2
)

//PushData
type PushData struct {
    Mid  string
    Ops  int32
    Data []byte
}

//ChanList
type ChanList struct {
    PushChan       chan *PushData
    CloseChan      chan string
    DisconnectChan chan string
}

type Comet struct {
    Sndbuf       int
    Rcvbuf       int
    KeepAlive    bool
    Reader       int
    ReadBuf      int
    ReadBufSize  int
    Writer       int
    WriteBuf     int
    WriteBufSize int
}

// TCP is tcp config.
type TCP struct {
    Bind []string
}

// Websocket is websocket config.
type Websocket struct {
    Bind        []string
    TLSOpen     bool
    TLSBind     []string
    CertFile    string
    PrivateFile string
}

// Protocol is proto config.
type Protocol struct {
    Proxy            bool
    Timer            int
    TimerSize        int
    SvrProto         int
    CliProto         int
    HandshakeTimeout xtime.Duration
    HeartbeatTimeout xtime.Duration
}

type Auth struct {
    Open  bool
    AppID string
}

// Bucket is bucket config.
type Bucket struct {
    Size    int
    Channel int
}

//ChanSize
type ChanSize struct {
    Push       int
    Close      int
    Disconnect int
}

// Config is comet config.
type ServerConfig struct {
    Comet     *Comet
    TCP       *TCP
    Websocket *Websocket
    Protocol  *Protocol
    Auth      *Auth
    Bucket    *Bucket
    ChanSize  *ChanSize
}
type methodHandler func(srv interface{}, ctx context.Context, data []byte, interceptor UnaryServerInterceptor) ([]byte, error)

type MethodDesc struct {
    Ops        int32
    MethodName string
    Handler    methodHandler
}

type ServiceDesc struct {
    ServiceName string
    // The pointer to the service interface. Used to check whether the user
    // provided implementation satisfies the interface requirements.
    HandlerType interface{}
    Methods     []MethodDesc
}

type service struct {
    server interface{}
    md     map[int32]*MethodDesc
}

// Server is comet server.
type Server struct {
    c              *ServerConfig
    round          *round.Round     // accept round store
    buckets        []*bucket.Bucket // subkey bucket
    bucketIdx      uint32
    serverID       string
    handlers       []UnaryServerInterceptor
    m              *service
    pushChan       chan *PushData
    closeChan      chan string
    disconnectChan chan string
    atuhClient     proto.AuthClient
}

func newAuthClient(appID string) (proto.AuthClient, error) {
    client := warden.DefaultClient()
    cc, err := client.Dial(context.Background(), fmt.Sprintf("discovery://default/%s", appID))
    if err != nil {
        return nil, err
    }
    return proto.NewAuthClient(cc), nil
}

func NewServer(conf *ServerConfig) (s *Server) {
    if conf == nil {
        conf = _defaultSerConf
    }
    roundConfig := round.RoundOptions{
        Reader:       conf.Comet.Reader,
        ReadBuf:      conf.Comet.ReadBuf,
        ReadBufSize:  conf.Comet.ReadBufSize,
        Writer:       conf.Comet.Writer,
        WriteBuf:     conf.Comet.WriteBuf,
        WriteBufSize: conf.Comet.WriteBufSize,
        Timer:        conf.Protocol.Timer,
        TimerSize:    conf.Protocol.TimerSize,
    }
    s = &Server{
        c:     conf,
        round: round.NewRound(roundConfig),
    }
    if conf.Auth.Open {
        var err error
        s.atuhClient, err = newAuthClient(conf.Auth.AppID)
        if err != nil {
            panic(err)
        }
    }
    // init bucket
    s.buckets = make([]*bucket.Bucket, conf.Bucket.Size)
    s.bucketIdx = uint32(conf.Bucket.Size)
    for i := 0; i < conf.Bucket.Size; i++ {
        s.buckets[i] = bucket.NewBucket(conf.Bucket.Channel)
    }
    s.pushChan = make(chan *PushData, conf.ChanSize.Push)
    s.closeChan = make(chan string, conf.ChanSize.Close)
    s.disconnectChan = make(chan string, conf.ChanSize.Disconnect)
    s.Use(s.recovery(), serverLogging(0))
    return
}

// Bucket get the bucket by subkey.
func (s *Server) GetBucket(subKey string) *bucket.Bucket {
    idx := cityhash.CityHash32([]byte(subKey), uint32(len(subKey))) % s.bucketIdx
    return s.buckets[idx]
}

func (s *Server) Use(handlers ...UnaryServerInterceptor) *Server {
    finalSize := len(s.handlers) + len(handlers)
    if finalSize >= int(_abortIndex) {
        panic("comet: server use too many handlers")
    }
    mergedHandlers := make([]UnaryServerInterceptor, finalSize)
    copy(mergedHandlers, s.handlers)
    copy(mergedHandlers[len(s.handlers):], handlers)
    s.handlers = mergedHandlers
    return s
}

func (s *Server) RegisterService(sd *ServiceDesc, ss interface{}) (cl *ChanList) {
    ht := reflect.TypeOf(sd.HandlerType).Elem()
    st := reflect.TypeOf(ss)
    if !st.Implements(ht) {
        log.Fatalf("comet: Server.RegisterService found the handler of type %v that does not satisfy %v", st, ht)
    }
    if s.m != nil {
        log.Fatalf("grpc: Server.RegisterService found duplicate service registration for %q", sd.ServiceName)
    }
    srv := &service{
        server: ss,
        md:     make(map[int32]*MethodDesc),
    }
    for i := range sd.Methods {
        d := &sd.Methods[i]
        srv.md[d.Ops] = d
    }
    s.m = srv
    cl = &ChanList{
        PushChan:       s.pushChan,
        CloseChan:      s.closeChan,
        DisconnectChan: s.disconnectChan,
    }
    //close
    go func() {
        for {
            mid := <-s.closeChan
            if channel := s.GetBucket(mid).Channel(mid); channel != nil {
                channel.Close()
            }
        }
    }()
    //push
    go func() {
        for {
            pd := <-s.pushChan
            s.PushByChannel(pd.Mid, pd.Ops, pd.Data)
        }
    }()
    return
}

func (s *Server) PushByChannel(sessionID string, ops int32, data []byte) {
    if channel := s.GetBucket(sessionID).Channel(sessionID); channel != nil {
        pushBody := &proto.Body{
            Ops:  ops,
            Data: data,
        }
        var (
            data []byte
            err  error
        )
        if data, err = gproto.Marshal(pushBody); err != nil {
            log.Println("push proto Marshal err:", err)
            return
        }
        p := &proto.Payload{
            Place: 1,
            Type:  int32(proto.Push),
            Body:  data,
        }
        if err := channel.Push(p); err != nil {
            log.Println("channle push err：", err)
        }
    }
}

func (s *Server) Operate(ctx context.Context, p *proto.Payload) (err error) {
    if p.Type == int32(proto.Request) {
        p.Type = int32(proto.Response)
    }
    reqBody := &proto.Body{}
    if err = gproto.Unmarshal(p.Body, reqBody); err != nil {
        return
    }
    srv := s.m
    md, ok := srv.md[reqBody.Ops]
    if !ok {
        return
    }
    reply, errCode := md.Handler(srv.server, ctx, reqBody.Data, s.interceptor)
    p.Code = int32(ecode.Cause(errCode).Code())
    p.Body = reply
    return
}

func (s *Server) interceptor(ctx context.Context, req interface{}, args *UnaryServerInfo, handler UnaryHandler) ([]byte, error) {
    var (
        i     int
        chain UnaryHandler
    )
    n := len(s.handlers)
    if n == 0 {
        return handler(ctx, req)
    }
    chain = func(ic context.Context, ir interface{}) ([]byte, error) {
        if i == n-1 {
            return handler(ic, ir)
        }
        i++
        return s.handlers[i](ic, ir, args, chain)
    }
    return s.handlers[0](ctx, req, args, chain)
}

func (s *Server) recoveryServer() (err error) {
    if rerr := recover(); rerr != nil {
        const size = 64 << 10
        buf := make([]byte, size)
        rs := runtime.Stack(buf, false)
        if rs > size {
            rs = size
        }
        buf = buf[:rs]
        pl := fmt.Sprintf("panic: %v\n%s\n", rerr, buf)
        fmt.Fprintf(os.Stderr, pl)
        err = fmt.Errorf(pl)
    }
    return
}

func (s *Server) Close() (err error) {
    return
}
