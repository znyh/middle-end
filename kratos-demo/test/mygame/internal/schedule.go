package internal

import (
    "flag"

    //"truco/internal/base/lib/hall"
    //"truco/internal/base/lib/rconfig"
    //"truco/internal/config"
    //"truco/internal/core/glog"
    //"truco/internal/core/playermgr"
    //"truco/internal/core/tablemgr"
    //testclient "truco/internal/core/test/client"
    //
    //"git.huoys.com/middle-end/kratos/pkg/log"

    "kratos-demo/test/mygame/internal/tablemgr"

    "github.com/go-kratos/kratos/v2/log"
)

var (
    style = 0
)

func init() {
    flag.IntVar(&style, "style", 0, "specify the style 1 - client; others - server.")
}

func Init() {
    switch style {
    case 1:
        client()
    default:
        server()
    }
}

func server() {

    log.Info("start server v0.0.1")

    //config.LoadConfig()
    //
    //rconfig.InitParams(&rconfig.RoomConfig{
    //	GameID:   int32(config.GameID),
    //	ArenaID:  int32(config.ArenaID),
    //	ServerID: config.ServerID,
    //}, config.GetTC(), config.GetGC(), config.GetCC(), config.GetRbC())
    //
    //hall.Init()

    //glog.OnStart()

    //playermgr.Init()
    //
    tablemgr.Init()
}

func client() {
    //testclient.OnStart()
}
