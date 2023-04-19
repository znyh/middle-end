package gomsg

import (
	"log"
	"sync/atomic"
	"time"
)

// Callback request callback
type Callback func(*Result)

// internal with keep alive handling
type iinternal interface {
	IHandler
	keepAlive()
}

type RequestHandler func(*Session, []byte) *Result
type PushHandler func(*Session, []byte) int16
type OpenHandler func(*Session)
type CloseHandler func(*Session, bool)

// IHandler node handler
type IHandler interface {
	OnOpen(*Session)
	OnClose(*Session, bool)
	OnReq(*Session, []byte) *Result
	OnPush(*Session, []byte) int16
}

// Node struct
type Node struct {
	addr            string
	internalHandler iinternal
	signal          chan int
	keepAliveSignal chan int
	ReadCounter     uint32
	WriteCounter    uint32
}

// NewNode make node ptr
func newNode(addr string, h iinternal) Node {
	return Node{
		addr:            addr,
		internalHandler: h,
		ReadCounter:     uint32(0),
		WriteCounter:    uint32(0),
		signal:          make(chan int),
		keepAliveSignal: make(chan int, 1),
	}
}

// Stop stop the IOCounter service
func (n *Node) Stop() {
	n.signal <- 1
	n.keepAliveSignal <- 1
	GetLoop().Stop()
}

// Start
func (n *Node) Start() {
	go n.ioCounter()
	go n.internalHandler.keepAlive()
	GetLoop().Start()
}

// IOCounter io couter
func (n *Node) ioCounter() {
	defer Recover()

	log.Println("IOCounter running.")
	d := time.Second * 5
	t := time.NewTimer(d)

	stop := false
	for !stop {
		select {

		case <-t.C:
			t.Reset(d)
			log.Printf("Read : %d/s\tWrite : %d/s\n", n.ReadCounter/5, n.WriteCounter/5)
			atomic.StoreUint32(&n.ReadCounter, 0)
			atomic.StoreUint32(&n.WriteCounter, 0)

		case <-n.signal:
			stop = true
		}
	}

	log.Println("IOCounter stoped.")
}
