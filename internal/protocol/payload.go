package protocol

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	PING byte = iota
	PONG
	AUTH
	SYNC
	SUSR
	SPWD
	SUCC
	FAIL

	MaxPayloadSize uint16 = 50 * (2 << 9) // 50KiB
)

var ErrMaxSizeExceeded = errors.New("maximum paylod size exceeded")

type Payload struct {
	payloadType byte
	bytes       []byte
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
	case AUTH:
		return "AUTH"
	case SYNC:
		return "SYNC"
	case SUSR:
		return "SUSR"
	case SPWD:
		return "SPWD"
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
	bytes := []byte{m.payloadType}
	lenBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(lenBytes, uint16(len(m.Bytes())))
	bytes = append(bytes, lenBytes...)
	bytes = append(bytes, m.Bytes()...)

	n, err := w.Write(bytes)
	if err != nil {
		return int64(n), err
	}

	return int64(n), nil
}

func (m *Payload) ReadFrom(r io.Reader) (int64, error) {
	err := binary.Read(r, binary.BigEndian, &m.payloadType)
	if err != nil {
		return 0, err
	}
	var n int64 = 1

	var size uint16
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

func NewSucc() *Payload {
	return &Payload{SUCC, []byte{}}
}

func NewSuccWithData(data []byte) *Payload {
	// Recepient will just have to know what to do with the data
	return &Payload{SUCC, data}
}

func NewFail(message string) *Payload {
	return &Payload{FAIL, []byte(message)}
}

func NewPayload(payloadType byte, bytes []byte) (*Payload, error) {
	if len(bytes) > int(MaxPayloadSize) {
		return nil, ErrMaxSizeExceeded
	}
	return &Payload{payloadType, bytes}, nil
}
