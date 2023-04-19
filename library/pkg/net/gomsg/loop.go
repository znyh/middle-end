package gomsg

//
//import (
//	"log"
//	"sync"
//)
//
//var (
//	once sync.Once
//	ins  *Loop
//)
//
//type Loop struct {
//	jobs   chan func()
//	toggle chan struct{}
//}
//
//func GetLoop() *Loop {
//	once.Do(func() {
//		ins = &Loop{
//			jobs: make(chan func(), 10000),
//		}
//	})
//
//	return ins
//}
//
//// Run starts the loop goroutine
//func (loop *Loop) Start() {
//	// the loop roroutine
//	go func() {
//		defer RecoverFromError(func() { loop.Start() })
//
//		for {
//			select {
//			case <-loop.toggle:
//				log.Println("Loop down.")
//				return
//			case job := <-loop.jobs:
//				job()
//			}
//		}
//	}()
//}
//
//// Run stop the loop goroutine
//func (loop *Loop) Stop() {
//	// prevent called from loop routine, so 'go' it here
//	go func() { loop.toggle <- struct{}{} }()
//}
//
//// PostPush a handler to excute in loop goroutine
//func (loop *Loop) Post(job func()) {
//	go func() { loop.jobs <- job }()
//}

import (
	"sync"
	"sync/atomic"
)

var (
	once sync.Once
	ins  *Loop

	bucketNum = 8
	bucketCap = 1024
)

type Loop struct {
	routines      []chan func()
	routinesCount uint64
}

func GetLoop() *Loop {
	once.Do(func() {
		ins = &Loop{}

		ins.routines = make([]chan func(), bucketNum)
		for i := 0; i < bucketNum; i++ {
			c := make(chan func(), bucketCap)
			ins.routines[i] = c
			go ins.startWorker(c)
		}
	})

	return ins
}

func (loop *Loop) startWorker(jobs chan func()) {

	for {
		select {
		case job, ok := <-jobs:
			if ok {
				job()
			}
		}
	}
}

func (loop *Loop) Start() {
	defer RecoverFromError(func() { loop.Start() })
}

func (loop *Loop) Post(job func()) {
	atomic.CompareAndSwapUint64(&loop.routinesCount, 1<<64-1, 0)
	num := int(atomic.AddUint64(&loop.routinesCount, 1)) % bucketNum
	loop.routines[num] <- job
}

func (loop *Loop) Stop() {
	for i := 0; i < int(bucketNum); i++ {
		close(loop.routines[i])
	}
}
