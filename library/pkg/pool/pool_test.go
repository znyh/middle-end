package pool

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestPost(t *testing.T) {
	l := NewPool()
	for j := 0; j < 5; j++ {
		l.Post(func() {
			myRand()
		})
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}

func TestPostAndWait(t *testing.T) {
	l := NewPool()
	for j := 0; j < 5; j++ {
		l.PostAndWait(func() interface{} {
			myRand()
			return nil
		})
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
