package proto

const (
	// OpProtoReady proto ready
	OpProtoReady = int32(1)
	// OpProtoFinish proto finish
	OpProtoFinish = int32(2)
)

type Pattern byte

const (
	Push Pattern = iota
	Request
	Response
	Ping
	Pong
	Sub
	Unsub
	Pub
)

const (
	AuthOps = -1
)
