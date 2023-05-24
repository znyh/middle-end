package gtable

import (
    "container/list"
    "fmt"
    "time"

    "kratos-demo/test/mygame/internal/base/rconfig/rconfig"
    "kratos-demo/test/mygame/internal/config"

    "github.com/go-kratos/kratos/v2/log"
)

const (
    TableStateStopped = 0
    TableStateRunning = 100
)

type Table struct {
    ID          int32  // 桌子ID
    MaxCnt      int16  // 最大玩家数
    sitCnt      int16  // 座位上的玩家
    uiGameState uint16 // 游戏状态
    IsClosed    bool   // 是否停服

    chairList []int32    // 玩家列表
    standList *list.List // 站立列表

    mLogic GameLogic // 桌子逻辑
}

//Init 初始化
func (tb *Table) Init() {
    tb.sitCnt = 0
    tb.chairList = make([]int32, tb.MaxCnt)
    tb.standList = list.New()
    tb.uiGameState = TableStateRunning //TableStateRunning //TableStateStopped
    tb.IsClosed = false

    tb.mLogic.init(tb)
}

//OnTimer 桌子定时
func (tb *Table) OnTimer() {
    tb.mLogic.run()
}

func (tb *Table) Empty() bool {
    return tb.sitCnt <= 0
}

//IsRunning run flag
func (tb *Table) IsRunning() bool {
    return TableStateRunning == tb.uiGameState
}

//StopTable 桌子停止事件
func (tb *Table) StopTable() {
    tb.uiGameState = TableStateStopped
    log.Infof("[table] stop table.")
}

func (tb *Table) start() {

    tb.uiGameState = TableStateRunning

    //日志
    tb.mLogic.mLog.CheckCreate(fmt.Sprintf(config.TABLE_LOG_PATH, rconfig.ArenaIDStr(),
        rconfig.ServerIDStr(), time.Now().Month(), time.Now().Day(), tb.ID))

    tb.mLogic.start()
}
