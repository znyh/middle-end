package gtimer

import "time"

/*
	游戏定时器
	注意：该定时器非线程安全，回调函数会在Update函数中执行。
*/

// 毫秒定时器
type stTimer struct {
	reg      int64
	period   int64
	f        func()
	repeated bool
}

// GameTimer 桌子定时器
type GameTimer struct {
	mapTimers map[int32]*stTimer
}

var gTimer GameTimer

func init() {
	gTimer.Init()
}

//Update global update
func Update() {
	gTimer.Update()
}

//Register register 毫秒级别
func Register(id int32, period int64, start bool, repeated bool, f func()) {
	gTimer.Register(id, period, start, repeated, f)
}

// Init 导出初始化
func (t *GameTimer) Init() {
	t.mapTimers = make(map[int32]*stTimer)
}

// Update 更新
func (t *GameTimer) Update() {
	if len(t.mapTimers) <= 0 {
		return
	}

	nowTick := getTick()

	for id, v := range t.mapTimers {
		if nowTick >= (v.reg + v.period) {
			v.f()
			if !v.repeated {
				delete(t.mapTimers, id)
			} else {
				v.reg = nowTick
			}
		}
	}

}

// Register 注册
func (t *GameTimer) Register(id int32, period int64, start bool, repeated bool, f func()) {
	if period <= 0 || f == nil {
		return
	}

	// set new timer
	t.mapTimers[id] = &stTimer{
		reg:      getTick(),
		period:   period,
		f:        f,
		repeated: repeated,
	}

	if start {
		f()
	}

}

// Exist 定时器存在
func (t *GameTimer) Exist(id int32) bool {
	_, ok := t.mapTimers[id]
	return ok
}

func (t *GameTimer) Cancel(id int32) {
	delete(t.mapTimers, id)
}

func getTick() int64 {
	return time.Now().UnixNano() / 1e6
}
