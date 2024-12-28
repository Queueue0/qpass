package protocol

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	PING byte = iota
	PONG
	SYNC
	SUCC
	FAIL

	MaxPayloadSize uint32 = 32 << 20
)

var ErrMaxSizeExceeded = errors.New("maximum paylod size exceeded")

type Payload struct {
	payloadType byte
	bytes []byte
}

func (m *Payload) Type() byte {
	return m.payloadType
}

// Returns a string representation of the payload type
// for logging and debugging purposes
func (m *Payload) TypeString() string {
	switch m.payloadType {
	case PING:
		return "PING"
	case PONG:
		return "PONG"
	case SYNC:
		return "SYNC"
	case SUCC:
		return "SUCC"
	case FAIL:
		return "FAIL"
	}

	return "INVALID TYPE"
}

func (m *Payload) Bytes() []byte {
	return m.bytes
}

func (m *Payload) String() string {
	return string(m.bytes)
}

func (m *Payload) WriteTo(w io.Writer) (int64, error) {
	err := binary.Write(w, binary.BigEndian, m.payloadType)
	if err != nil {
		return 0, err
	}
	var n int64 = 1

	err = binary.Write(w, binary.BigEndian, uint32(len(m.bytes)))
	if err != nil {
		return n, err
	}
	n += 4

	b, err := w.Write(m.bytes)
	if err != nil {
		return n, err
	}

	return n + int64(b), nil
}

func (m *Payload) ReadFrom(r io.Reader) (int64, error) {
	err := binary.Read(r, binary.BigEndian, &m.payloadType)
	if err != nil {
		return 0, err
	}
	var n int64 = 1

	var size uint32
	err = binary.Read(r, binary.BigEndian, &size)
	if err != nil {
		return 0, err
	}
	n += 4
	
	if size > MaxPayloadSize {
		return n, ErrMaxSizeExceeded
	}

	m.bytes = make([]byte, size)
	b, err := r.Read(m.bytes)

	if err != nil {
		return n, err
	}

	return n + int64(b), err
}

func NewPing() *Payload {
	return &Payload{PING, []byte{}}
}

func NewPong() *Payload {
	return &Payload{PONG, []byte{}}
}

func NewPayload(payloadType byte, bytes []byte) (*Payload, error) {
	if uint32(len(bytes)) > MaxPayloadSize {
		return nil, ErrMaxSizeExceeded
	} 
	return &Payload{payloadType, bytes}, nil
}
