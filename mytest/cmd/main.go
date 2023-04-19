package main

import (
    "flag"
    "net"
    "os"
    "os/signal"
    "syscall"
    "time"

    "mytest/internal/conf"

    "github.com/gin-gonic/gin"
    "github.com/go-kratos/kratos/v2/config"
    "github.com/go-kratos/kratos/v2/config/file"
    "github.com/go-kratos/kratos/v2/log"
    "google.golang.org/grpc"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
    // Name is the name of the compiled software.
    Name string = "hello.service"
    // Version is the version of the compiled software.
    Version string = "v1"
    // flagconf is the config flag.
    flagconf string

    id, _ = os.Hostname()
)

func init() {
    flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

type (
    App struct {
        hs *gin.Engine
        gs *grpc.Server
    }
)

func newApp(hs *gin.Engine, gs *grpc.Server) *App {
    return &App{
        hs: hs,
        gs: gs,
    }
}

func main() {
    flag.Parse()
    c := config.New(
        config.WithSource(
            file.NewSource(flagconf),
        ),
    )
    defer c.Close()

    if err := c.Load(); err != nil {
        panic(err)
    }

    var bc conf.Bootstrap
    if err := c.Scan(&bc); err != nil {
        panic(err)
    }

    app, closeFunc, err := wireApp(bc.Server, bc.Data)
    if err != nil {
        closeFunc()
        panic(err)
    }
    app.run(&bc)

    ch := make(chan os.Signal, 1)
    signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
    for {
        s := <-ch
        log.Info("get a signal %s", s.String())
        switch s {
        case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
            closeFunc()
            log.Info("account exit")
            time.Sleep(time.Second)
            return
        case syscall.SIGHUP:
        default:
            return
        }
    }
}

func (a *App) run(c *conf.Bootstrap) {

    go func() {
        // 监听一个端口
        listen, err := net.Listen("tcp", c.Server.Grpc.Addr)
        if err != nil {
            log.Fatal(err)
        }
        // 启动服务
        if err = a.gs.Serve(listen); err != nil {
            log.Fatal(err)
        }
    }()

    go func() {
        log.Fatal(a.hs.Run(c.Server.Http.Addr))
    }()
}
