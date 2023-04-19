package calc

import (
    "log"
    "runtime/debug"
    "time"
)

const (
    _PoolCap = 10240
)

//var (
//   _PoolCap = 10240
//)

type Pool struct {
    jobs chan func()
}

func NewPool() *Pool {
    p := &Pool{
        jobs: make(chan func(), _PoolCap),
    }
    p.start()
    return p
}

func (p *Pool) start() {
    go func() {
        defer p.Recover(func() {
            p.start()
        })
        for {
            select {
            case job, ok := <-p.jobs:
                if ok {
                    job()
                }
            }
        }
    }()
}

// 异步投递线程函数
func (p *Pool) Post(job func()) {
    go func() {
        p.jobs <- job
    }()
}

// 同步投递线程函数
func (p *Pool) PostAndWait(job func() interface{}) interface{} {
    ch := make(chan interface{})
    go func() {
        p.jobs <- func() {
            ch <- job()
        }
    }()
    return <-ch
}

func (p *Pool) Recover(cb func()) {
    if e := recover(); e != nil {
        timestamp := time.Now().Format("2006-01-02 15:04:05")
        log.Printf("[%+v] Recover => %s:%s\n", timestamp, e, debug.Stack())

        if cb != nil {
            cb()
        }
    }
}
