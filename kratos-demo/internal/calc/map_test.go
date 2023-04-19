package calc

import (
    "fmt"
    "sync/atomic"
    "testing"
    "time"

    "kratos-demo/internal/base"

    "github.com/go-kratos/kratos/v2/log"
)

type buck struct {
    routines []chan int

    routinesNum uint64

    RoutineAmount uint64

    m *Map
}

func newBuck() *buck {
    b := new(buck)

    b.RoutineAmount = 8
    b.routinesNum = 0

    b.routines = make([]chan int, b.RoutineAmount)

    for i := uint64(0); i < b.RoutineAmount; i++ {
        c := make(chan int, 1024)
        b.routines[i] = c
        go b.PushRoom(c)
    }

    b.m = New(10)

    return b
}

func (b *buck) PushRoom(c chan int) {
    //fmt.Printf("创建pushroom\n")
    for {
        arg := <-c

        //fmt.Println("arg:", arg)

        if arg < 0 {

        }

        b.m.calc()

        fmt.Println("arg:", arg)

        //if room = b.Room(arg.RoomId); room != nil {
        //	room.Push(&arg.P)
        //}

    }

}

func (b *buck) BroadcastRoom(arg int) {
    // 广播消息递增id
    num := atomic.AddUint64(&b.routinesNum, 1) % b.RoutineAmount
    fmt.Printf("BroadcastRoom RoomMsgArg :%d bucket routinesNum :%d\n", arg, b.routinesNum)
    b.routines[num] <- arg
}

//func main2(betCnt int) {
//   ///newBuck()
//   b := newBuck()
//   b.BroadcastRoom(1234)
//
//   b.BroadcastRoom(1235)
//
//}

func TestNew(t *testing.T) {

    var (
        betCnt = 3
        digCnt = 10
    )

    {
        start := time.Now()

        b := newBuck()
        for i := 0; i < betCnt; i++ {
            b.BroadcastRoom(i)
        }
        defer log.Infof("use:%+v/ms _totalCnt:%+v _winCnt(%+v) _loseCnt(%+v)\n", time.Since(start).Milliseconds(), b.m._totalCnt, b.m._winCnt, b.m._loseCnt)
    }

    {
        var (
            _winCnt   = int64(0)
            _loseCnt  = int64(0)
            _totalCnt = int64(0)
        )
        start := time.Now()
        for i := 0; i < betCnt; i++ {
            m := New(digCnt)

            m.calc()

            atomic.AddInt64(&_totalCnt, 1)
            if base.IsHit(30) {
                atomic.AddInt64(&_winCnt, 1)
            } else {
                atomic.AddInt64(&_loseCnt, 1)
            }
        }
        log.Infof("use:%+v/ms _totalCnt:%+v _winCnt(%+v) _loseCnt(%+v)\n", time.Since(start).Milliseconds(), _totalCnt, _winCnt, _loseCnt)
    }

    //pool := NewPool()
    //{
    //    var (
    //        _winCnt   = int64(0)
    //        _loseCnt  = int64(0)
    //        _totalCnt = int64(0)
    //    )
    //
    //    start := time.Now()
    //    for i := 0; i < betCnt; i++ {
    //
    //        pool.Post(func() {
    //            m := New(digCnt)
    //
    //            m.calc()
    //
    //            atomic.AddInt64(&_totalCnt, 1)
    //            if base.IsHit(30) {
    //                atomic.AddInt64(&_winCnt, 1)
    //            } else {
    //                atomic.AddInt64(&_loseCnt, 1)
    //            }
    //        })
    //
    //    }
    //
    //    log.Infof("use:%+v/ms _totalCnt:%+v _winCnt(%+v) _loseCnt(%+v)\n", time.Since(start).Milliseconds(), _totalCnt, _winCnt, _loseCnt)
    //}
    //
    //{
    //    var (
    //        _winCnt   = int64(0)
    //        _loseCnt  = int64(0)
    //        _totalCnt = int64(0)
    //    )
    //    start := time.Now()
    //    for i := 0; i < betCnt; i++ {
    //
    //        pool.PostAndWait(func() interface{} {
    //            m := New(digCnt)
    //
    //            m.calc()
    //
    //            atomic.AddInt64(&_totalCnt, 1)
    //            if base.IsHit(30) {
    //                atomic.AddInt64(&_winCnt, 1)
    //            } else {
    //                atomic.AddInt64(&_loseCnt, 1)
    //            }
    //            return nil
    //
    //        })
    //    }
    //    log.Infof("use:%+v/ms _totalCnt:%+v _winCnt(%+v) _loseCnt(%+v)\n", time.Since(start).Milliseconds(), _totalCnt, _winCnt, _loseCnt)
    //}

}

//
//func doProducer(x int, digCnt int, out chan<- int) {
//    for {
//
//        m := New(digCnt)
//
//        val := m.calc()
//
//        out <- val
//    }
//}

//
//func main21() {
//    ch := make(chan int, 10)
//    go producer(2, ch)
//    go producer(3, ch)
//    go consumer(ch)
//    sig := make(chan os.Signal)
//    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
//    <-sig
//}
//
//func producer(x int, out chan<- int) {
//    var i = 0
//    for {
//        out <- i * x
//        i++
//    }
//}
//
//func consumer(in <-chan int) {
//    for i := range in {
//        log.Printf("i = %d", i)
//    }
//}
