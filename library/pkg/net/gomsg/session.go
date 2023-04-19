package gomsg

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// Result (ErrorNum,Data)
type Result struct {
	En   int16
	Data []byte
}

// Succeed if ret is success
func (ret *Result) Succeed() bool {
	return ret.En == 0
}

// Session 会话
type Session struct {
	ID   int32
	conn net.Conn
	node *Node

	bodyLen      uint16
	reqSeed      uint32   //seq
	reqPool      sync.Map //map[uint32]chan *Result    // seq->result
	closed       bool
	responseTime int64
}

// NewSession make session
func newSession(id int32, conn net.Conn, n *Node) *Session {
	return &Session{
		ID:           id,
		conn:         conn,
		node:         n,
		reqPool:      sync.Map{},
		responseTime: time.Now().Unix(),
	}
}

func (s *Session) split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	len := len(data)
	offset := 0
	if len == 0 {
		return 0, nil, nil
	}

	if atEOF {
		return len, nil, nil
	}

	if s.bodyLen == 0 {
		if len < 2 {
			// Request more data.
			return 0, nil, nil
		}

		s.bodyLen = binary.LittleEndian.Uint16(data[offset:2])
		len -= 2
		offset += 2

		if len < int(s.bodyLen) {
			return 2, nil, nil
		}

	} else if len < int(s.bodyLen) {
		// Request more data.
		return 0, nil, nil
	}

	advance = int(s.bodyLen) + offset
	s.bodyLen = 0
	return advance, data[offset:advance], nil
}

func (s *Session) scan() {
	defer Recover()

	input := bufio.NewScanner(s.conn)
	input.Split(s.split)

	for input.Scan() {
		// dispatch
		s.dispatch(input.Bytes())

		atomic.AddUint32(&s.node.ReadCounter, 1)
		//s.node.ReadCounter <- 1
	}

	s.Close(false)
}

func (s *Session) elapsedSinceLastResponse() int {
	return int(time.Now().Unix() - atomic.LoadInt64(&s.responseTime))
}

func (s *Session) dispatch(data []byte) {
	// log.Printf("conn : %d=> Read [% x]\n", s.ID, data)
	atomic.StoreInt64(&s.responseTime, time.Now().Unix())

	if len(data) < 2 {
		return
	}

	reader := bytes.NewBuffer(data)
	_, err := reader.ReadByte() // cnt
	if err != nil {
		s.Close(false)
		log.Println(err.Error())
		return
	}

	pattern, err := reader.ReadByte()
	if err != nil {
		s.Close(false)
		log.Println(err.Error())
		return
	}

	left := len(data) - 2

	switch Pattern(pattern) {
	case Push:
		s.onPush(reader, left)
	case Request:
		s.onReq(reader, left)
	case Response:
		s.onResponse(reader, left)
	case Ping:
		s.onPing(reader)
	case Pong:
		s.onPong(reader)
	case Sub:
	case Unsub:
	case Pub:
	}
}

func (s *Session) onPush(reader *bytes.Buffer, left int) {
	body := make([]byte, left)
	n, err := reader.Read(body)
	if n != left || err != nil {
		s.Close(false)
		log.Println("")
		return
	}

	GetLoop().Post(func() {
		ret := s.node.internalHandler.OnPush(s, body)
		if ret != 0 {
			log.Printf("onPush : %d\n", ret)
		}
	})
}

func (s *Session) onReq(reader *bytes.Buffer, left int) {
	var serial uint16
	err := binary.Read(reader, binary.LittleEndian, &serial)
	if err != nil {
		s.Close(false)
		log.Println(err.Error())
		return
	}

	left -= 2
	body := make([]byte, left)
	n, err := reader.Read(body)
	if n != left {
		s.Close(false)
		log.Println("")
		return
	}

	if err != nil {
		s.Close(false)
		log.Println(err.Error())
		return
	}

	GetLoop().Post(func() {
		r := s.node.internalHandler.OnReq(s, body)
		if r == nil {
			log.Println("OnReq returns nil!!!")
			return
		}

		err = s.response(serial, r)
		if err != nil {
			log.Println(err)
		}
	})
}

func (s *Session) onResponse(reader *bytes.Buffer, left int) {
	var serial uint16
	var en int16
	err := binary.Read(reader, binary.LittleEndian, &serial)
	if err != nil {
		s.Close(false)
		log.Println(err.Error())
		return
	}

	left -= 2
	err = binary.Read(reader, binary.LittleEndian, &en)
	if err != nil {
		s.Close(false)
		log.Println(err.Error())
		return
	}

	left -= 2
	body := make([]byte, left)
	n, err := reader.Read(body)
	if n != left {
		s.Close(false)
		log.Println("")
		return
	}

	if err != nil {
		s.Close(false)
		log.Println(err.Error())
		return
	}

	req, exists := s.reqPool.Load(serial)
	if !exists {
		log.Printf("%d not exist in req pool.\n", serial)
		return
	}

	s.reqPool.Delete(serial)
	req.(chan *Result) <- &Result{En: en, Data: body}
}

func (s *Session) onPing(reader *bytes.Buffer) {
	serial, err := reader.ReadByte()
	if err != nil {
		s.Close(false)
		log.Println(err.Error())
		return
	}

	err = s.pong(serial)
	if err != nil {
		log.Println(err)
	}
}

// todo lizs
func (s *Session) onPong(reader *bytes.Buffer) {
}

// pong   //14seq...
func (s *Session) pong(serial byte) error {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, uint16(1+1+1))
	if err != nil {
		return err
	}
	buf.WriteByte(1)
	buf.WriteByte(byte(Pong))
	buf.WriteByte(serial)

	err = s.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// response // 12xxx...
func (s *Session) response(serial uint16, ret *Result) error {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, uint16(1+1+2+2+len(ret.Data)))
	if err != nil {
		return err
	}
	buf.WriteByte(1)
	buf.WriteByte(byte(Response))
	err = binary.Write(buf, binary.LittleEndian, serial)
	if err != nil {
		return err
	}
	err = binary.Write(buf, binary.LittleEndian, ret.En)
	if err != nil {
		return err
	}
	if len(ret.Data) != 0 {
		buf.Write(ret.Data)
	}

	err = s.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// Close 关闭会话
func (s *Session) Close(force bool) {
	GetLoop().Post(func() {
		if s.closed {
			return
		}

		s.closed = true
		s.conn.Close()

		log.Printf("conn [%d] closed.\n", s.ID)
		s.node.internalHandler.OnClose(s, force)
	})
}

// IsClosed 是否关闭
func (s *Session) Closed() bool {
	return s == nil || s.closed
}

// Ping ping remote, returns delay seconds  //ping:130...
func (s *Session) Ping() error {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, uint16(1+1+1))
	if err != nil {
		return err
	}
	buf.WriteByte(1)
	buf.WriteByte(byte(Ping))
	buf.WriteByte(0)

	err = s.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// Write raw send interface
func (s *Session) Write(data []byte) error {
	// log.Printf("conn : %d=> Write [% x]\n", s.ID, data)
	n, err := s.conn.Write(data)
	if n != len(data) {
		if err != nil {
			log.Println("Write error =>", err.Error())
			s.Close(false)
		} else {
			log.Printf("Write error => writed : %d != expected : %d", n, len(data))
		}
	} else {
		//s.node.WriteCounter <- 1
		atomic.AddUint32(&s.node.WriteCounter, 1)
	}

	return err
}

// Request request remote to response  //10==xxx...
func (s *Session) Request(data []byte) *Result {
	var en NetError

	for {
		if len(data) == 0 {
			en = RequestDataIsEmpty
			break
		}

		serial := uint16(atomic.AddUint32(&s.reqSeed, 1))
		if _, exists := s.reqPool.Load(serial); exists {
			en = SerialConflict
			break
		}

		req := make(chan *Result)
		// record, let 'response' package know which chan to notify
		s.reqPool.Store(serial, req)

		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, uint16(1+1+2+len(data)))
		if err != nil {
			en = Write
			break
		}
		buf.WriteByte(1)
		buf.WriteByte(byte(Request))
		err = binary.Write(buf, binary.LittleEndian, serial)
		if err != nil {
			en = Write
			break
		}
		if len(data) != 0 {
			_, err = buf.Write(data)
			if err != nil {
				en = Write
				break
			}
		}

		err = s.Write(buf.Bytes())
		if err != nil {
			en = Write
			break
		}

		return <-req
	}

	return &Result{En: int16(en)}
}

// Push push to remote without response  // 10xxx...
func (s *Session) Push(data []byte) NetError {
	if len(data) == 0 {
		return PushDataIsEmpty
	}

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, uint16(1+1+len(data)))
	if err != nil {
		return Write
	}

	buf.WriteByte(1)
	buf.WriteByte(byte(Push))
	if len(data) != 0 {
		_, err = buf.Write(data)
		if err != nil {
			return Write
		}
	}

	err = s.Write(buf.Bytes())
	if err != nil {
		return Write
	}

	return Success
}

// Pub ...   //17xx0xxx...
func (s *Session) Pub(subject string, data []byte) error {
	subBytes := []byte(subject)

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, uint16(1+1+len(subBytes)+1+len(data)))
	if err != nil {
		return err
	}
	buf.WriteByte(1)
	buf.WriteByte(byte(Pub))
	buf.Write(subBytes)
	buf.WriteByte(0) // append \0 to string
	buf.Write(data)

	err = s.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// Sub ...  //15xxx...
func (s *Session) Sub(subject string) error {
	subBytes := []byte(subject)

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, uint16(1+1+len(subBytes)))
	if err != nil {
		return err
	}
	buf.WriteByte(1)
	buf.WriteByte(byte(Sub))
	buf.Write(subBytes)

	err = s.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) GetRemoteIP() string {
	return s.conn.RemoteAddr().String()
}
