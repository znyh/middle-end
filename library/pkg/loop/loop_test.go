package loop

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestPost(t *testing.T) {
	l := NewLoop()
	start := time.Now()
	for j := 0; j < 3; j++ {
		l.Post(func() {
			myRand()
		})
	}
	fmt.Printf("since:%+v\n", time.Since(start))
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}

func TestPostAndWait(t *testing.T) {
	l := NewLoop()
	start := time.Now()
	for j := 0; j < 3; j++ {
		l.PostAndWait(func() interface{} {
			myRand()
			return nil
		})
	}
	fmt.Printf("since:%+v\n", time.Since(start))
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}

func TestPostAndWait2(t *testing.T) {

	ch := make(chan func(), 1024)

	go func() {
		for {
			select {
			case job, ok := <-ch:
				if ok {
					job()
				}
			}
		}
	}()

	for j := 0; j < 3; j++ {
		ch <- func() { myRand() }
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}

var i = 0

func myRand() {
	i++
	s := time.Duration(i) * time.Second
	time.Sleep(s)
	fmt.Printf("=====> i:%d timerNow:%+v\n", i, time.Now().Second())
}
