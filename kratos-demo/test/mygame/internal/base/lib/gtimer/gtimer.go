package gtimer

import (
    "time"

    "kratos-demo/test/mygame/internal/base/lib"
)

/*
	全局定时器
*/

func Once(duration time.Duration, f func()) {
    run(duration, duration, false, f)
}

func Forever(duration time.Duration, f func()) {
    run(duration, duration, true, f)
}

func ForeverNow(duration time.Duration, f func()) {
    lib.Post(f)
    Forever(duration, f)
}

func ForeverTime(durFirst, durRepeat time.Duration, f func()) {
    run(durFirst, durRepeat, true, f)
}

func run(durFirst, durRepeat time.Duration, repeated bool, f func()) {
    go func() {
        timer := time.NewTimer(durFirst)
        for {
            select {
            case <-timer.C:
                lib.Post(f)
                if repeated {
                    timer.Reset(durRepeat)
                } else {
                    return
                }
            }
        }
    }()
}
