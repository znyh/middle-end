package gtable

import (
    "kratos-demo/test/mygame/internal/base"
    "kratos-demo/test/mygame/internal/config"
    "kratos-demo/test/mygame/internal/glog"
)

type GameLogic struct {
    mTable *Table
    mLog   glog.TableLog

    nTimeOut int64 //超时
    nStage   int32 //阶段

    activeChair int32 //当前操作玩家
    bankerChair int32 //庄家玩家
}

//初始化桌子
func (tl *GameLogic) init(pTable *Table) {
    tl.mTable = pTable
    tl.mLog.Init(pTable.ID)
}

//桌子启动
func (tl *GameLogic) start() {
    tl.updateStage(config.StPrepare)
}

func (tl *GameLogic) updateStage(stage int32) {

    tl.mLog.Stage(tl.nStage, stage) //状态转移日志

    tl.nStage = stage

    tl.nTimeOut = config.GetStageInter(stage)

    tl.checkResetTimer()
}

//设置超时时间
func (tl *GameLogic) setTimer(n int32) {
    tl.nTimeOut = base.GetTick() + int64(n*1000)
}

func (tl *GameLogic) checkResetTimer() {
}
