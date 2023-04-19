package proto

import (
	"encoding/binary"
	"errors"

	"github.com/znyh/middle-end/library/pkg/net/comet/internal/bufio"
	"github.com/znyh/middle-end/library/pkg/net/comet/internal/websocket"
)

const (
    // MaxBodySize max proto body size
    MaxBodySize = int32(1 << 12)
)

const (
    //size
    _packSize      = 2
    _placeSize     = 1
    _typeSize      = 1
    _seqSize       = 2
    _heartSeqSize  = 1
    _codeSize      = 2
    _reqHeaderSize = _packSize + _placeSize + _typeSize
    _maxPackSize   = int32(_reqHeaderSize) + int32(_seqSize) + MaxBodySize
    // offset
    _packOffset  = 0
    _placeOffset = _packOffset + _packSize
    _typeOffset  = _placeOffset + _placeSize
    _seqOffset   = _typeOffset + _typeSize
    _codeOffset  = _seqOffset + _codeSize
)

var (
    // ErrProtoPackLen proto packet len error
    ErrProtoPackLen = errors.New("default server codec pack length error")
)

var (
    // ProtoReady proto ready
    ProtoReady = &Payload{Op: OpProtoReady}
    // ProtoFinish proto finish
    ProtoFinish = &Payload{Op: OpProtoFinish}
)

func (p *Payload) ReadTCP(rr *bufio.Reader) (err error) {
    var (
        bodyLen int
        packLen int32
        buf     []byte
    )
    if buf, err = rr.Pop(_reqHeaderSize); err != nil {
        return
    }
    packLen = int32(binary.LittleEndian.Uint16(buf[_packOffset:_placeOffset]))
    p.Place = int32(buf[_placeOffset])
    p.Type = int32(buf[_typeOffset])
    if packLen > _maxPackSize {
        return ErrProtoPackLen
    }
    bodyLen = int(packLen - int32(_placeSize) - int32(_typeSize))
    if bodyLen < 1 {
        return
    }
    if buf, err = rr.Pop(bodyLen); err != nil {
        return
    }
    if p.Type == int32(Ping) || p.Type == int32(Pong) {
        p.Seq = int32(buf[0])
        p.Body = nil
    } else if p.Type == int32(Push) {
        p.Body = buf
    } else if p.Type == int32(Request) && bodyLen > _seqSize {
        p.Seq = int32(binary.LittleEndian.Uint16(buf[0:]))
        p.Body = buf[_seqSize:]
    } else if p.Type == int32(Response) && bodyLen > _reqHeaderSize {
        p.Seq = int32(binary.LittleEndian.Uint16(buf[0:]))
        p.Code = int32(binary.LittleEndian.Uint16(buf[_seqSize:]))
        p.Body = buf[_seqSize+_codeSize:]
    }
    return
}

func (p *Payload) WriteTCP(wr *bufio.Writer) (err error) {
    var (
        buf        []byte
        packLen    int
        headerSize int
    )
    if Pattern(p.Type) == Response {
        headerSize = _placeSize + _typeSize + _seqSize + _codeSize
    } else if Pattern(p.Type) == Request {
        headerSize = _placeSize + _typeSize + _seqSize
    } else {
        headerSize = _placeSize + _typeSize
    }
    if len(p.Body) > int(MaxBodySize) {
        return ErrProtoPackLen
    }
    packLen = headerSize + len(p.Body)
    if buf, err = wr.Peek(headerSize + _packSize); err != nil {
        return
    }
    binary.LittleEndian.PutUint16(buf[_packOffset:], uint16(packLen))
    buf[_placeOffset] = byte(1)
    buf[_typeOffset] = byte(p.Type)
    if Pattern(p.Type) == Response {
        binary.LittleEndian.PutUint16(buf[_seqOffset:], uint16(p.Seq))
        binary.LittleEndian.PutUint16(buf[_codeOffset:], uint16(p.Code))
    } else if Pattern(p.Type) == Request {
        binary.LittleEndian.PutUint16(buf[_seqOffset:], uint16(p.Seq))
    }
    if p.Body != nil {
        _, err = wr.Write(p.Body)
    }
    return
}

func (p *Payload) WriteTCPHeart(wr *bufio.Writer) (err error) {
    var (
        buf     []byte
        packLen int
    )
    packLen = _placeSize + _typeSize + _heartSeqSize
    dataLen := _packSize + packLen
    if buf, err = wr.Peek(dataLen); err != nil {
        return
    }
    // header
    binary.LittleEndian.PutUint16(buf[_packOffset:], uint16(packLen))
    buf[_placeOffset] = byte(1)
    buf[_typeOffset] = byte(p.Type)
    buf[_seqOffset] = byte(p.Seq)
    return
}

// ReadWebsocket read a proto from websocket connection.
func (p *Payload) ReadWebsocket(ws *websocket.Conn) (err error) {
    var (
        buf []byte
    )
    if _, buf, err = ws.ReadMessage(); err != nil {
        return
    }

    dataLen := len(buf)
    if dataLen < (_reqHeaderSize - _heartSeqSize) {
        return ErrProtoPackLen
    }
    p.Place = int32(buf[0])
    p.Type = int32(buf[_placeSize])
    seqPos := _placeSize + _typeSize
    if p.Type == int32(Ping) {
        p.Seq = int32(buf[seqPos])
        p.Body = nil
    } else if p.Type == int32(Push) {
        p.Body = buf[seqPos:]
    } else if dataLen > _reqHeaderSize {
        p.Seq = int32(binary.LittleEndian.Uint16(buf[seqPos:]))
        pos := seqPos + _seqSize
        p.Body = buf[pos:]
    }
    return
}

// WriteWebsocket write a proto to websocket connection.
func (p *Payload) WriteWebsocket(ws *websocket.Conn) (err error) {
    var (
        buf        []byte
        packLen    int
        headerSize int
    )
    if Pattern(p.Type) == Response {
        headerSize = _placeSize + _typeSize + _seqSize + _codeSize
    } else {
        headerSize = _placeSize + _typeSize
    }
    packLen = headerSize + len(p.Body)
    if err = ws.WriteHeader(websocket.BinaryMessage, packLen); err != nil {
        return
    }
    if buf, err = ws.Peek(headerSize); err != nil {
        return
    }
    buf[0] = byte(1)
    buf[_placeSize] = byte(p.Type)
    if Pattern(p.Type) == Response {
        pos := _placeSize + _typeSize
        binary.LittleEndian.PutUint16(buf[pos:], uint16(p.Seq))
        pos += _seqSize
        binary.LittleEndian.PutUint16(buf[pos:], uint16(p.Code))
    }
    if p.Body != nil {
        err = ws.WriteBody(p.Body)
    }
    return
}

// WriteWebsocketHeart write websocket heartbeat with room online.
func (p *Payload) WriteWebsocketHeart(wr *websocket.Conn) (err error) {
    var (
        buf     []byte
        packLen int
    )
    packLen = _placeSize + _typeSize + _heartSeqSize
    // websocket header
    if err = wr.WriteHeader(websocket.BinaryMessage, packLen); err != nil {
        return
    }
    if buf, err = wr.Peek(packLen); err != nil {
        return
    }
    // header
    buf[0] = byte(1)
    pos := _placeSize
    buf[pos] = byte(Pong)
    pos += _typeSize
    buf[pos] = byte(p.Seq)
    return
}
