package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	"gomsg"
)

type handler struct {
}

func (h *handler) request(s *gomsg.Session) {
	ret := s.Request([]byte("hello"))
	if ret.En == int16(gomsg.Success) {
		go h.request(s)
	}
}

func (h *handler) OnOpen(s *gomsg.Session) {

	fmt.Println("handler OnOpen...")
	go h.request(s)
}

func (h *handler) OnClose(s *gomsg.Session, force bool) {
}

func (h *handler) OnReq(s *gomsg.Session, data []byte) *gomsg.Result {
	fmt.Println("handler onReq --> ", string(data))
	return &gomsg.Result{En: int16(gomsg.Success), Data: nil}
}

func (h *handler) OnPush(s *gomsg.Session, data []byte) int16 {
	return int16(gomsg.Success)
}

type shandler struct {
}

func (h *shandler) OnOpen(s *gomsg.Session) {
	fmt.Println("shandler OnOpen...")
}

func (h *shandler) OnClose(s *gomsg.Session, force bool) {
}

func (h *shandler) OnReq(s *gomsg.Session, data []byte) *gomsg.Result {
	//fmt.Println("shandler onReq --> ", string(data))
	return &gomsg.Result{En: int16(gomsg.Success), Data: nil}
}

func (h *shandler) OnPush(s *gomsg.Session, data []byte) int16 {
	return int16(gomsg.Success)
}

func main() {
	host := flag.String("h", "localhost:6000", "specify the client/server host address.\n\tUsage: -h localhost:6000")
	runAsServer := flag.Bool("s", true, "whether to run as a tcp server.\n\tUsage : -s true/false")
	parallel := flag.Int("p", 1, "Parallel count.\n\tUsage : -p 32")
	flag.Parse()

	go func() {
		log.Println(http.ListenAndServe("localhost:8080", nil))
	}()

	if *runAsServer {
		s := gomsg.NewServer(*host, &shandler{})
		go s.Start()
	} else {
		for i := 0; i < *parallel; i++ {
			c := gomsg.NewClient(*host, &handler{}, true)
			c.Start()
		}
	}

	ch := make(chan int)
	<-ch
}
