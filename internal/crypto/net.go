package crypto

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/binary"
	"net"
	"slices"
	"time"

	"golang.org/x/crypto/argon2"
)

func genSharedKey(b []byte) []byte {
	return argon2.IDKey(b, nil, 3, 64*1024, 2, 32)
}

const pubKeySize = 32

type secureConn struct {
	c  net.Conn
	ss []byte // Shared Secret, generated by dh
}

/*
	 This is currently just an ECDH key exchange
	 TODO:
		Server long-term keys
		Server Signature
		Server MAC
*/
func NewClientConn(c net.Conn) (*secureConn, error) {
	privkey, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	pubkey := privkey.PublicKey()
	_, err = c.Write(pubkey.Bytes())
	if err != nil {
		return nil, err
	}

	b := make([]byte, pubKeySize)
	_, err = c.Read(b)
	if err != nil {
		return nil, err
	}

	remoteKey, err := ecdh.X25519().NewPublicKey(b)
	if err != nil {
		return nil, err
	}

	ss, err := privkey.ECDH(remoteKey)
	if err != nil {
		return nil, err
	}

	ss = genSharedKey(ss)

	return &secureConn{c, ss}, nil
}

func Dial(addr string) (*secureConn, error) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return NewClientConn(c)
}

func NewServerConn(c net.Conn) (*secureConn, error) {
	b := make([]byte, pubKeySize)
	_, err := c.Read(b)
	if err != nil {
		return nil, err
	}

	remoteKey, err := ecdh.X25519().NewPublicKey(b)
	if err != nil {
		return nil, err
	}

	privkey, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	ss, err := privkey.ECDH(remoteKey)
	if err != nil {
		return nil, err
	}

	ss = genSharedKey(ss)

	pubkey := privkey.PublicKey()
	_, err = c.Write(pubkey.Bytes())
	if err != nil {
		return nil, err
	}

	return &secureConn{c, ss}, nil
}

func (s *secureConn) Read(b []byte) (int, error) {
	// Size schenanigans might be unnecessary, added them while debugging
	// I think it probably makes it more robust anyway, so I'm leaving it for now
	sizeBytes := make([]byte, 4)
	_, err := s.c.Read(sizeBytes)
	if err != nil {
		return 0, err
	}

	size := binary.BigEndian.Uint32(sizeBytes)

	buf := make([]byte, size)
	_, err = s.c.Read(buf)
	if err != nil {
		return 0, err
	}

	d, err := DecryptBytes(buf, s.ss)
	if err != nil {
		return 0, err
	}

	n := min(len(b), len(d))

	for i := 0; i < n; i++ {
		b[i] = d[i]
	}

	return n, nil
}

func (s *secureConn) Write(b []byte) (int, error) {
	e, err := EncryptBytes(b, s.ss)
	if err != nil {
		return 0, err
	}

	sizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBytes, uint32(len(e)))

	n, err := s.c.Write(slices.Concat(sizeBytes, e))
	if err != nil {
		if n > len(b) {
			n = len(b)
		}
		return n, err
	}

	return len(b), nil
}

func (s *secureConn) Close() error {
	return s.c.Close()
}

func (s *secureConn) LocalAddr() net.Addr {
	return s.c.LocalAddr()
}

func (s *secureConn) RemoteAddr() net.Addr {
	return s.c.RemoteAddr()
}

func (s *secureConn) SetDeadline(t time.Time) error {
	return s.c.SetDeadline(t)
}

func (s *secureConn) SetReadDeadline(t time.Time) error {
	return s.c.SetReadDeadline(t)
}

func (s *secureConn) SetWriteDeadline(t time.Time) error {
	return s.c.SetWriteDeadline(t)
}
