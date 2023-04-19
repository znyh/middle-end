package loop

import (
	"log"
	"runtime/debug"
	"time"
)

const (
	_LoopCap = 10240
)

type Loop struct {
	jobs chan func()
}

func NewLoop() *Loop {
	p := &Loop{
		jobs: make(chan func(), _LoopCap),
	}
	p.start()
	return p
}

func (p *Loop) start() {
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
func (p *Loop) Post(job func()) {
	go func() {
		p.jobs <- job
	}()
}

// 同步投递线程函数
func (p *Loop) PostAndWait(job func() interface{}) interface{} {
	ch := make(chan interface{})
	go func() {
		p.jobs <- func() {
			ch <- job()
		}
	}()
	return <-ch
}

func (p *Loop) Recover(cb func()) {
	if e := recover(); e != nil {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		log.Printf("[%+v] Recover => %s:%s\n", timestamp, e, debug.Stack())

		if cb != nil {
			cb()
		}
	}
}
