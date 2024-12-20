package protocol

import (
	"encoding/binary"
	"net"
)

func Read(c net.Conn) []byte {
	l := make([]byte, 2)
	i, _ := c.Read(l)

	lm := binary.BigEndian.Uint16(l[:i])

	b := make([]byte, lm)
	e, _ := c.Read(b)

	return b[:e]
}

func Write(t byte, m []byte, c net.Conn) {
	p := make([]byte, 1)
	p[0] = t
	l := make([]byte, 2)
	binary.BigEndian.PutUint16(l, uint16(len(m)))
	p = append(p, l...)
	c.Write(append(p, m...))
}
