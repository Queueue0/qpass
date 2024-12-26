package protocol

import (
	"encoding/binary"
	"net"
)

func Read(c net.Conn) Packet {
	l := make([]byte, 2)
	i, _ := c.Read(l)

	lm := binary.BigEndian.Uint16(l[:i])

	b := make([]byte, lm)
	e, _ := c.Read(b)

	switch b[0] {
	case PING:
		return NewPing()
	case PONG:
		return NewPong()
	case SYNC:
		lm := binary.BigEndian.Uint16(b[1:i])
		
		d := make([]byte, )
	}

	// TODO: process other types of packets


	return NewPong()
}

func Write(c net.Conn, p Packet) error {
	pb := make([]byte, 1)
	pb[0] = p.Type()
	l := make([]byte, 2)
	m := p.Data()
	binary.BigEndian.PutUint16(l, uint16(len(m)))
	pb = append(pb, l...)
	_, err := c.Write(append(pb, m...))
	if err != nil {
		return err
	}

	return nil
}
