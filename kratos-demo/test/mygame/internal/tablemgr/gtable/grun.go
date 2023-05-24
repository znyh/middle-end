package gtable

import (
    "kratos-demo/test/mygame/internal/base"
    "kratos-demo/test/mygame/internal/config"

    "github.com/go-kratos/kratos/v2/log"
)

// 桌子运行时的状态,整个行牌流程
func (tl *GameLogic) run() {
    // 超时处理
    if base.GetTick() > tl.nTimeOut {
        switch tl.nStage {
        case config.StPrepare:
            tl.prepareStart()
        case config.StSendCard:
            tl.onSendCardTimeout()
        case config.StOutCard:

        case config.StLaunch:

        case config.StResult:

        default:
        }
        log.Infof("timeout ... %d %d  ....", tl.nStage, tl.activeChair)
    }
}

//发牌超时
func (tl *GameLogic) prepareStart() {
    tl.mLog.WriteLog("prepareStart timeout")
    tl.mTable.start()
    tl.mLog.WriteLog("prepareStart start")
    tl.mTable.StopTable()
}

//发牌超时
func (tl *GameLogic) onSendCardTimeout() {
    tl.mLog.WriteLog("onSendCardTimeout timeout")
}
