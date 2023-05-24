package tablemgr

import (
    "time"

    "kratos-demo/test/mygame/internal/base/lib/gtimer"
    "kratos-demo/test/mygame/internal/tablemgr/gtable"

    "github.com/go-kratos/kratos/v2/log"
)

var (
    mapTables   map[int32]*gtable.Table
    listTables  []*gtable.Table
    IsCloseRoom bool
)

//Init init
func Init() {

    tableNum := 2
    chairNum := 4

    mapTables = make(map[int32]*gtable.Table)
    listTables = make([]*gtable.Table, tableNum)
    IsCloseRoom = false

    for i := 1; i <= tableNum; i++ {
        table := &gtable.Table{
            ID:     int32(i),
            MaxCnt: int16(chairNum),
        }
        table.Init()
        table.ID = int32(i)
        mapTables[table.ID] = table
        listTables[i-1] = table
    }

    gtimer.Forever(time.Second/2, OnTimer)

    //// 初始化机器人
    //if config.GetRbC().IsOpen {
    //    robotmgr.Init()
    //}

    log.Infof("桌子初始化完成. 桌子数:%d 椅子:%d", len(mapTables), chairNum)
}

// 房间定时器
func OnTimer() {
    for _, v := range mapTables {
        if v.IsRunning() {
            v.OnTimer()
        }
    }

}
