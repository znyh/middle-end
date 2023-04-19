package gomsg

// NetError network error code
type NetError int16

const (
	Success NetError = iota
	ExceptionCatched
	Write
	Read
	RequestDataIsEmpty
	SerialConflict
	NoHandler
	ReadErrorNo
	SessionClosed
	PushDataIsEmpty
)

// Pattern msg pattern
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
