package base

const (
	DayInter = 86400 //DayInter one day inter second
)

//TimeZone time
func TimeZone() int64 {
	return 8
}

//GetWeets 0点时间戳
func GetWeets() int64 {
	now := GetTick() / 1000
	weeNowTs := now - (now+TimeZone()*3600)%DayInter
	return weeNowTs
}

//IsToday 判定时间戳是否为今天
func IsToday(tick int64) bool {
	now := GetTick() / 1000
	tick = tick / 1000
	weeNowTs := now - (now+TimeZone()*3600)%DayInter
	weeTick := tick - (tick+TimeZone()*3600)%DayInter

	return weeTick == weeNowTs
}
