package glog

import (
    "flag"
    "fmt"
    "log"
    "os"
    "time"

    //"truco/internal/base"
    //"truco/internal/base/lib/gtimer"
    //"truco/internal/base/lib/rconfig"

    "kratos-demo/test/mygame/internal/base"
    "kratos-demo/test/mygame/internal/base/lib/gtimer"
    "kratos-demo/test/mygame/internal/base/rconfig/rconfig"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var (
    Room     *zap.SugaredLogger
    Table    *zap.SugaredLogger
    createTs int64
    PATH     string
)

func init() {
    flag.StringVar(&PATH, "zap", os.Getenv("ZAP_PATH"), "specify log path.")
    if PATH == "" {
        PATH = "./log/truco/"
    }
}

//OnStart 启动日志
func OnStart() {
    log.Printf(" zaplog onstart ... ")
    //reload()
    gtimer.ForeverNow(5*60*time.Second, reload)
}

func reload() {
    if createTs > 0 && base.IsToday(createTs) {
        return
    }

    t := time.Now()
    timeStr := fmt.Sprintf("%04d-%02d-%02d_%02d_%02d_%02d",
        t.Year(), t.Month(), t.Day(),
        t.Hour(), t.Minute(), t.Second())

    Room = initLog(PATH, fmt.Sprintf("room_%s_%s_%s", rconfig.ArenaIDStr(), rconfig.ServerIDStr(), timeStr))
    Table = initLog(PATH, fmt.Sprintf("table_%s_%s_%s", rconfig.ArenaIDStr(), rconfig.ServerIDStr(), timeStr))

    createTs = base.GetTick()
}

func initLog(path, name string) *zap.SugaredLogger {
    if _, err := os.Stat(path); err != nil {
        os.MkdirAll(path, 0777)
    }

    cfg := zap.NewProductionConfig()
    cfg.EncoderConfig.TimeKey = "time"
    cfg.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
    cfg.OutputPaths = []string{
        fmt.Sprintf("%s/%s.log", path, name),
    }
    cfg.ErrorOutputPaths = []string{
        "stderr",
    }

    logger, err := cfg.Build()
    if err != nil {
        panic(err)
    }

    return logger.Named(name).Sugar()
}
