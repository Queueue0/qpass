package protocol

import (
	"encoding/binary"
	"net"
)

func Read(c net.Conn) Packet {
	t := make([]byte, 1)
	c.Read(t)

	switch t[0] {
	case PING:
		return NewPing()
	case PONG:
		return NewPong()
	}

	// TODO: process other types of packets

	// l := make([]byte, 2)
	// i, _ := c.Read(l)

	// lm := binary.BigEndian.Uint16(l[:i])

	// b := make([]byte, lm)
	// e, _ := c.Read(b)

	return NewPing()
}

func Write(c net.Conn, p Packet) {
	pb := make([]byte, 1)
	pb[0] = p.Type()
	l := make([]byte, 2)
	m := p.Data()
	binary.BigEndian.PutUint16(l, uint16(len(m)))
	pb = append(pb, l...)
	c.Write(append(pb, m...))
}
