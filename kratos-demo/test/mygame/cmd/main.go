package main

import (
    "os"
    "os/signal"
    "syscall"
    "time"

    "kratos-demo/test/mygame/internal"

    "github.com/go-kratos/kratos/v2/log"
)

func main() {
    internal.Init()

    c := make(chan os.Signal, 1)
    signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
    for {
        s := <-c
        log.Info("get a signal %s", s.String())
        switch s {
        case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
            log.Info("truco exit")
            time.Sleep(time.Second)
            return
        case syscall.SIGHUP:
        default:
            return
        }
    }
}
