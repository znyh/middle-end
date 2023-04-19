package channel

import (
    "sync"

    "github.com/znyh/middle-end/library/pkg/net/comet/internal/bufio"
    "github.com/znyh/middle-end/library/pkg/net/comet/internal/ring"
    "github.com/znyh/middle-end/library/pkg/net/comet/proto"
)

// Channel used by message pusher send msg to write goroutine.
type Channel struct {
    CliProto ring.Ring
    signal   chan *proto.Payload
    Writer   bufio.Writer
    Reader   bufio.Reader
    Next     *Channel
    Prev     *Channel

    Mid   int64
    Key   string
    IP    string
    mutex sync.RWMutex
}

// NewChannel new a channel.
func NewChannel(cli, svr int) *Channel {
    c := new(Channel)
    c.CliProto.Init(cli)
    c.signal = make(chan *proto.Payload, svr)
    return c
}

// Push server push message.
func (c *Channel) Push(p *proto.Payload) (err error) {
    select {
    case c.signal <- p:
    default:
    }
    return
}

// Ready check the channel ready or close?
func (c *Channel) Ready() *proto.Payload {
    return <-c.signal
}

// Signal send signal to the channel, proto ready.
func (c *Channel) Signal() {
    c.signal <- proto.ProtoReady
}

// Close close the channel.
func (c *Channel) Close() {
    c.signal <- proto.ProtoFinish
}
