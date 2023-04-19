package gomsg

import (
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// Server struct
type Server struct {
	Node
	sessions    sync.Map
	sessionsCnt int32
	seed        int32 //session id
	handler     IHandler
}

// Count return sessions'count
func (s *Server) Count() int {
	return int(s.sessionsCnt)
}

// OnOpen ...
func (s *Server) OnOpen(session *Session) {
	s.sessions.Store(s.seed, session)
	atomic.AddInt32(&s.sessionsCnt, 1)
	log.Printf("conn [%d] established.\n", session.ID)

	go s.handler.OnOpen(session)
}

// OnClose ...
func (s *Server) OnClose(session *Session, force bool) {
	s.sessions.Delete(session.ID)
	atomic.AddInt32(&s.sessionsCnt, -1)

	go s.handler.OnClose(session, force)
}

// OnReq ...
func (s *Server) OnReq(session *Session, data []byte) *Result {
	defer Recover()

	return s.handler.OnReq(session, data)
}

// OnPush ...
func (s *Server) OnPush(session *Session, data []byte) int16 {
	return s.handler.OnPush(session, data)
}

// keep alive
func (s *Server) keepAlive() {
	defer Recover()

	d := time.Second * 5
	t := time.NewTimer(d)

	stop := false
	for !stop {
		select {
		case <-t.C:
			t.Reset(d)
			s.sessions.Range(func(k, v interface{}) bool {
				s := v.(*Session)
				if s.elapsedSinceLastResponse() > 60 {
					s.Close(true)
				} else if s.elapsedSinceLastResponse() > 40 {
					err := s.Ping()
					if err != nil {
						log.Println(err)
					}
				}
				return true
			})

		case <-s.keepAliveSignal:
			stop = true
		}
	}

	log.Println("Keep alive stopped.")
}

// NewServer new tcp server
func NewServer(host string, h IHandler) *Server {
	ret := &Server{
		sessions: sync.Map{},
		seed:     0,
		handler:  h,
	}

	ret.Node = newNode(host, ret)
	return ret
}

// Start server startup
func (s *Server) Start() {
	listener, err := net.Listen("tcp4", s.addr)
	if err != nil {
		log.Panicf("listen failed : %v", err)
	} else {
		log.Printf("server running on [%s]\n", s.addr)
	}

	defer listener.Close()

	// base.start
	s.Node.Start()

	// accept incomming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err.Error())
			continue
		}

		s.handleConn(conn)
	}
}

// Stop server shutdown
func (s *Server) Stop() {
	s.sessions.Range(func(k, v interface{}) bool {
		v.(*Session).Close(true)
		s.sessions.Delete(k)
		return true
	})
	atomic.StoreInt32(&s.sessionsCnt, 0)

	s.keepAliveSignal <- 1
	s.Node.Stop()
	log.Println("server stopped.")
}

func (s *Server) handleConn(conn net.Conn) {
	atomic.AddInt32(&s.seed, 1)

	// make session
	session := newSession(s.seed, conn, &s.Node)

	// notify
	s.OnOpen(session)

	// io
	go session.scan()
}
