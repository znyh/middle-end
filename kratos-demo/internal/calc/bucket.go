package calc

//
//import (
//    "log"
//    "math"
//    "runtime/debug"
//    "sync/atomic"
//    "time"
//)
//
//const (
//    _bucketNum = 100
//    _bucketCap = 102400
//)
//
//type Pool struct {
//    buckets []chan func()
//    count   uint64
//}
//
//func NewPool() *Pool {
//    ins := &Pool{}
//
//    ins.buckets = make([]chan func(), _bucketNum)
//    for i := 0; i < _bucketNum; i++ {
//        c := make(chan func(), _bucketCap)
//        ins.buckets[i] = c
//        ins.startWorker(c)
//    }
//
//    return ins
//}
//
//func (p *Pool) startWorker(jobs chan func()) {
//    go func(jobs chan func()) {
//
//        defer p.Recover(func() {
//            p.startWorker(jobs)
//        })
//
//        for {
//            select {
//            case job, ok := <-jobs:
//                if ok {
//                    job()
//                }
//            }
//        }
//    }(jobs)
//}
//
//func (p *Pool) Post(job func()) {
//    atomic.CompareAndSwapUint64(&p.count, math.MaxUint64, 0)
//    num := int(atomic.AddUint64(&p.count, 1)) % (_bucketNum - 1)
//    p.buckets[num] <- job
//}
//
//func (p *Pool) PostAndWait(job func() interface{}) interface{} {
//    ch := make(chan interface{})
//    go func() {
//        p.buckets[_bucketNum-1] <- func() {
//            ch <- job()
//        }
//    }()
//    return <-ch
//
//    //ch := make(chan interface{})
//    //
//    //atomic.CompareAndSwapUint64(&p.count, math.MaxUint64, 0)
//    //num := int(atomic.AddUint64(&p.count, 1)) % (_bucketNum - 1)
//    //
//    //go func() {
//    //    p.buckets[num] <- func() {
//    //        ch <- job()
//    //    }
//    //}()
//    //return <-ch
//
//}
//
//func (p *Pool) Stop() {
//    for i := 0; i < int(_bucketNum); i++ {
//        close(p.buckets[i])
//    }
//}
//
//func (p *Pool) Recover(cb func()) {
//
//    if e := recover(); e != nil {
//        timestamp := time.Now().Format("2006-01-02 15:04:05")
//        log.Printf("[%+v] Recover => %s:%s\n", timestamp, e, debug.Stack())
//
//        if cb != nil {
//            cb()
//        }
//    }
//}
