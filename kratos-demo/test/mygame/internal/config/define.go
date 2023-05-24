package config

const (
	TABLE_LOG_PATH = "log/Room_%s_%s/%d_%d/Table_%d.log"
)

//阶段定义
const (
	StPrepare  = 0 //准备期/空闲期
	StSendCard = 1 //发牌期
	StOutCard  = 2 //出牌
	StLaunch   = 3 //加倍
	StResult   = 4 //结算
)

//阶段名称
var (
	stageDesc = map[int32]string{
		StPrepare:  "准备期阶段",
		StSendCard: "发牌期",
		StOutCard:  "等待出牌",
		StLaunch:   "加倍",
		StResult:   "结算",
	}
)

//阶段名
func StageDesc(status int32) string {
	return stageDesc[status]
}

func IsLaiZi(v int32) bool {
	return false
}
