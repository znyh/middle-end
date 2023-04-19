package gomsg

import (
	"log"
	"net"
	"time"
)

// Client struct
type Client struct {
	Node
	session          *Session
	autoRetryEnabled bool
	handler          IHandler
	//sta              *STAService
}

func (c *Client) OnOpen(s *Session) {
	c.session.responseTime = time.Now().Unix()
	c.handler.OnOpen(s)
}

func (c *Client) OnClose(s *Session, force bool) {
	c.handler.OnClose(s, force)

	// reconnect
	if !force && c.autoRetryEnabled {
		log.Println("reconnecting ...")
		c.Stop()
		c.Start()
	}
}

func (c *Client) OnReq(s *Session, data []byte) *Result {
	return c.handler.OnReq(s, data)
}

func (c *Client) OnPush(s *Session, data []byte) int16 {
	return c.handler.OnPush(s, data)
}

// keep alive
func (c *Client) keepAlive() {
	defer Recover()

	log.Println("Keep alive running.")

	d := time.Second * 5
	t := time.NewTimer(d)

	stop := false
	for !stop {
		select {
		case <-t.C:
			t.Reset(d)
			if nil == c.session {
				break
			}

			if c.session.elapsedSinceLastResponse() > 30 {
				log.Println("session not reponse over 30s, closing...")
				c.session.Close(false)
			} else if c.session.elapsedSinceLastResponse() > 10 {
				err := c.session.Ping()
				if err != nil {
					log.Println(err)
				}
			}

		case <-c.keepAliveSignal:
			stop = true
		}
	}

	log.Println("Keep alive stopped.")
}

// NewClient new tcp client
func NewClient(host string, h IHandler, autoRetry bool) *Client {
	ret := &Client{
		autoRetryEnabled: autoRetry,
		session:          nil,
		handler:          h,
	}
	ret.Node = newNode(host, ret)

	return ret
}

// Start client startup
func (c *Client) Start() {
	var conn net.Conn
	var err error
	for {
		conn, err = net.Dial("tcp4", c.Node.addr)
		if err != nil {
			log.Printf("connect failed : %v\n", err)

			if c.autoRetryEnabled {
				select {
				case <-time.After(time.Second * 2):
					log.Printf("reconnecting ...")
				}

				continue
			}

			log.Println("you can set `autoRetryEnabled` true to do auto reconnect stuff.")
			return
		}

		break
	}

	// base start
	c.Node.Start()

	// make session
	c.session = newSession(0, conn, &c.Node)

	// notify
	c.OnOpen(c.session)

	// io
	go c.session.scan()

	log.Printf("conn [%d] established.\n", c.session.ID)
}

// Stop client shutdown
func (c *Client) Stop() {
	c.Node.Stop()

	if c.session != nil {
		c.session.Close(true)
		c.session = nil
	}

	log.Println("client stopped.")
}

func (c *Client) Push(data []byte) NetError {
	return c.session.Push(data)
}

func (c *Client) Request(data []byte) *Result {
	return c.session.Request(data)
}
