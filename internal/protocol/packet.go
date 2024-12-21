package protocol

const (
	PING byte = iota
	PONG
	SYNC
	SUCC
	FAIL
)

type Packet interface {
	Type() byte
	Length() uint16
	Data() []byte
	ToBytes() []byte
}

type PingPong struct {
	ptype byte
}

func (p PingPong) Type() byte {
	return p.ptype
}

func (p PingPong) Length() uint16 {
	return uint16(0)
}

func (p PingPong) Data() []byte {
	return []byte{}
}

func (p PingPong) ToBytes() []byte {
	return []byte{p.ptype, 0, 0}
}

func NewPing() PingPong {
	return PingPong{PING}
}

func NewPong() PingPong {
	return PingPong{PONG}
}
