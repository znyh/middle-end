package config

import (
    "time"

    "kratos-demo/test/mygame/internal/base"
)

/*全局变量,不能声明桌子特殊属性变量 如:桌子规则 不同桌子规则不一样*/

var (
    CurTableNum = int32(0)
)

// 设置当前桌子数
func SetCurTableNum(tableNum int32) {
    CurTableNum = tableNum
}

// 获取当前桌子数
func GetCurTableNum() int32 {
    return CurTableNum
}

//NeedKick 需要踢出玩家
func NeedKickMin(m int64) bool {
    return m < GetGC().MinMoney
}
func NeedKickMax(m int64) bool {
    if GetGC().MaxMoney == -1 {
        return false
    }
    return m > GetGC().MaxMoney
}

func BaseScore() int64 {
    return GetGC().Coin.BaseScore
}

func TaiFlee() int64 {
    return GetGC().Coin.TaiFlee
}

func AvoidHurt() float32 {
    return GetGC().Coin.AvoidHurt
}

// 获取房间最大人数
func GetMaxPlayer() int32 {
    return int32(GetTC().TableNum * GetTC().ChairNum)
}

func NeedTakeIn() bool {
    return false
}

//获取场机器人思考时间区间
func GetAIThinkingTimeRange() time.Duration {
    start, end := 1, 5
    if ttr := GetRbC().ThinkTimeRange; len(ttr) == 2 {
        start = int(ttr[0])
        end = int(ttr[1])
    }
    return time.Second * time.Duration(base.RandRange(start, end))
}
