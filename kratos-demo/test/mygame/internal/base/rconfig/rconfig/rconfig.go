package rconfig

import (
    "fmt"
    "strconv"

    "github.com/go-kratos/kratos/v2/log"
)

var (
    //Room 房间信息
    ins *RoomConfig = &RoomConfig{
        GameID:   12001,
        ArenaID:  1,
        ServerID: "truco",
    }
)

//RoomConfig 入场配置
type RoomConfig struct {
    GameID   int32
    ArenaID  int32
    ServerID string
}

func GameID() int32 {
    return ins.GameID
}

func ArenaID() int32 {
    return ins.ArenaID
}

func ArenaIDStr() string {
    return strconv.Itoa(int(ins.ArenaID))
}

func ServerIDStr() string {
    return ins.ServerID
}

//InitParams 场配置初始化
func InitParams(r *RoomConfig, c ...interface{}) {
    ins = r

    str := ""
    for _, v := range c {
        str += fmt.Sprintf("%+v,  ", v)
    }
    log.Info("\r\n\t**房间基本信息**: %+v \r\n\tIns:%+v", ins, str)
}
