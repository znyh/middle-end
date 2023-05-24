package glog

import (
    "fmt"
    "log"
    "os"
    "strings"
    "time"

    //"truco/internal/base"
    //"truco/internal/base/lib/rconfig"
    //"truco/internal/config"
    //"truco/internal/core/playermgr/gplayer"
    //"truco/internal/core/settle"

    "kratos-demo/test/mygame/internal/base"
    "kratos-demo/test/mygame/internal/base/rconfig/rconfig"
    "kratos-demo/test/mygame/internal/config"
)

/*
	桌子日志
*/

//TableLog .
type TableLog struct {
    tableID  int32
    initFlag bool
    nLogTs   int64 // 日志创建时间
    gameNo   string

    _file  *os.File
    logger *log.Logger
}

func (tlog *TableLog) Init(tableID int32) {
    tlog.nLogTs = 0
    tlog.initFlag = false
    tlog.tableID = tableID

    //log
    filePath := fmt.Sprintf(config.TABLE_LOG_PATH, rconfig.ArenaIDStr(),
        rconfig.ServerIDStr(), time.Now().Month(), time.Now().Day(), tableID)
    tlog.CheckCreate(filePath)

}

func (tlog *TableLog) Create(fileName string) {
    //if !config.GetGC().LocalLog {
    //    return
    //}

    idx := strings.LastIndex(fileName, "/")
    path := fileName[0:idx]
    if _, err := os.Stat(path); err != nil {
        os.MkdirAll(path, 0777)
    }

    tlog.Close()

    _file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
    if err != nil {
        return
    }
    tlog.logger = log.New(_file, "", log.Ltime)
    tlog.initFlag = true
    tlog.nLogTs = base.GetTick()
}

func (tlog *TableLog) Close() {
    if tlog._file != nil {
        tlog._file.Close()
    }
}

func (tlog *TableLog) CheckCreate(fileName string) {
    if base.IsToday(tlog.nLogTs) {
        return
    }

    tlog.Close()
    tlog.Create(fileName)
}

func (tlog *TableLog) WriteLog(msg string, args ...interface{}) {
    content := fmt.Sprintf(msg, args...)

    //Table.Infow("玩家操作",
    //    "场ID", rconfig.ArenaIDStr(),
    //    "房间ID", rconfig.ServerIDStr(),
    //    "桌子ID", tlog.tableID,
    //    "局号", tlog.gameNo,
    //    "内容", content)

    if tlog.initFlag == false {
        return
    }

    tlog.logger.Printf(content)
}

//游戏开始
func (tlog *TableLog) GameBegin(nRound int32, gameNo string) {
    tlog.gameNo = gameNo

    logs := []string{fmt.Sprintf("\n\t一局开始，当前局:%d", nRound)}

    //for _, p := range players {
    //    if p != nil {
    //        logs = append(logs, fmt.Sprintf("玩家[%d, %d] 金币[%d] ",
    //            p.GetPlayerID(), p.GetChairID(), p.GetMoney()))
    //    }
    //}

    tlog.WriteLog(strings.Join(logs, "\r\n"))
}

//游戏结束
func (tlog *TableLog) GameEnd(result string) {
    logs := []string{"\n\t一局结算"}
    //obj.RangePlayer(func(k int32, p *gplayer.Player) bool {
    //    logs = append(logs, fmt.Sprintf("玩家[%d, %d] 输赢分[%d] 桌费[%d], 金币[%d]",
    //        p.GetPlayerID(), p.GetChairID(), p.GetWin(), p.PlayerBase.GameHallData.Tax, p.GetMoney()))
    //
    //    return true
    //})
    tlog.WriteLog(strings.Join(logs, result))
    tlog.WriteLog(strings.Join(logs, "\r\n"))
}

func (tlog *TableLog) Stage(old int32, new int32) {
    if old == new {
        return
    }
    tlog.WriteLog("【状态转移】%s -- > %s", config.StageDesc(old), config.StageDesc(new))
}
