package lib

import (
    "runtime/debug"

    "github.com/go-kratos/kratos/v2/log"
)

var (
    ins *Loop
)

func init() {
    ins = &Loop{
        jobs:   make(chan func(), 10240),
        toggle: make(chan byte, 1),
    }
    ins.Start()
}

type Loop struct {
    jobs   chan func()
    toggle chan byte
}

func RecoverFromError(cb func()) {
    if e := recover(); e != nil {
        log.Error("Recover => %s:%s\n", e, debug.Stack())
        if nil != cb {
            cb()
        }
    }
}
func (lp *Loop) Start() {
    log.Info("loop routine start.")
    go func() {
        defer RecoverFromError(func() {
            lp.Start()
        })
        for {
            select {
            case <-lp.toggle:
                log.Info("Loop routine stop.")
                return
            case job := <-lp.jobs:
                job()
            }
        }
    }()
}
func (lp *Loop) Stop() {
    go func() {
        lp.toggle <- 1
    }()
}
func Stop() { ins.Stop() }
func (lp *Loop) Jobs() int {
    return len(lp.jobs)
}
func Jobs() int { return ins.Jobs() }
func (lp *Loop) Post(job func()) {
    go func() {
        lp.jobs <- job
    }()
}
func Post(job func()) { ins.Post(job) }
func (lp *Loop) PostAndWait(job func() interface{}) interface{} {
    ch := make(chan interface{})
    go func() {
        lp.jobs <- func() {
            ch <- job()
        }
    }()
    return <-ch
}
func PostAndWait(job func() interface{}) interface{} { return ins.PostAndWait(job) }
