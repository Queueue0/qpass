package crypto

import (
	"bytes"
	"crypto"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/binary"
	"errors"
	"net"
	"slices"
	"time"

	"github.com/Queueue0/qpass/internal/structures"
	"golang.org/x/crypto/argon2"
	_ "golang.org/x/crypto/blake2b"
)

func genSharedKey(b []byte) []byte {
	return argon2.IDKey(b, nil, 3, 64*1024, 2, 32)
}

const (
	pubKeySize    = 32
	rsaKeyByteLen = 2
)

type secureConn struct {
	c     net.Conn
	ss    []byte                 // Shared Secret, generated by dh, hashed with argon2
	queue structures.Queue[byte] // For storing leftover bytes if the buffer supplied to Read isn't big enough
}

func NewClientConn(c net.Conn) (*secureConn, error) {
	// Generate ephemeral DH key pair
	privkey, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		c.Close()
		return nil, err
	}
	pubkey := privkey.PublicKey()

	// Send client hello
	// Initial packet to server containing our DH public key
	_, err = c.Write(pubkey.Bytes())
	if err != nil {
		c.Close()
		return nil, err
	}

	// Receive server's DH public key
	rkBytes := make([]byte, pubKeySize)
	_, err = c.Read(rkBytes)
	if err != nil {
		c.Close()
		return nil, err
	}

	// Receive server's RSA public key
	rsaKeyLenBuff := make([]byte, rsaKeyByteLen)
	_, err = c.Read(rsaKeyLenBuff)
	if err != nil {
		c.Close()
		return nil, err
	}

	rsaKeyLen := binary.BigEndian.Uint16(rsaKeyLenBuff)

	rsaKeyBytes := make([]byte, rsaKeyLen)
	_, err = c.Read(rsaKeyBytes)
	if err != nil {
		c.Close()
		return nil, err
	}

	rsaKey, err := x509.ParsePKCS1PublicKey(rsaKeyBytes)
	if err != nil {
		c.Close()
		return nil, err
	}

	knownHosts, err := getKnownHosts()
	if err != nil {
		c.Close()
		return nil, err
	}

	// Compare received RSA public key against our known hosts
	knownKey, ok := knownHosts[c.RemoteAddr().String()]
	if !ok {
		// Add host if it doesn't exist
		// TODO: probably add more checks for this case
		// Example: ssh will ask if you want to trust a new server
		addHost(c.RemoteAddr().String(), rsaKey)
	} else {
		// Verify that key matches the key we already have
		if !rsaKey.Equal(knownKey) {
			return nil, errors.New("Key does not match known key for this host")
		}
	}

	// Receive signature over all data exchanged so far
	sigLenBuff := make([]byte, 2)
	_, err = c.Read(sigLenBuff)
	if err != nil {
		c.Close()
		return nil, err
	}

	sigLen := binary.BigEndian.Uint16(sigLenBuff)

	sig := make([]byte, sigLen)
	_, err = c.Read(sig)
	if err != nil {
		c.Close()
		return nil, err
	}

	// Verify signature to confirm the server has the RSA private key
	opts := rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
		Hash:       crypto.BLAKE2b_512,
	}

	hash := opts.HashFunc().New()
	_, err = hash.Write(slices.Concat(pubkey.Bytes(), rkBytes, rsaKeyBytes))
	if err != nil {
		c.Close()
		return nil, err
	}

	err = rsa.VerifyPSS(rsaKey, crypto.BLAKE2b_512, hash.Sum(nil), sig, &opts)
	if err != nil {
		c.Close()
		return nil, err
	}

	// Compute DH shared secret
	remoteKey, err := ecdh.X25519().NewPublicKey(rkBytes)
	if err != nil {
		c.Close()
		return nil, err
	}

	ss, err := privkey.ECDH(remoteKey)
	if err != nil {
		c.Close()
		return nil, err
	}

	ss = genSharedKey(ss)

	// Receive MAC from server
	macLenBytes := make([]byte, 2)
	_, err = c.Read(macLenBytes)
	if err != nil {
		c.Close()
		return nil, err
	}

	macLen := binary.BigEndian.Uint16(macLenBytes)

	mac := make([]byte, macLen)
	_, err = c.Read(mac)
	if err != nil {
		c.Close()
		return nil, err
	}

	// Verify MAC to confirm server computed the same DH shared secret
	macData, err := decryptBytes(mac, ss, nil)
	if err != nil {
		c.Close()
		return nil, err
	}

	if !bytes.Equal(macData, slices.Concat(pubkey.Bytes(), rkBytes, rsaKeyBytes, sig)) {
		c.Close()
		return nil, errors.New("MAC authentication failed")
	}

	return &secureConn{c, ss, structures.Queue[byte]{}}, nil
}

// Just makes it easier to create a client-side secureConn
func Dial(addr string) (*secureConn, error) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return NewClientConn(c)
}

func NewServerConn(c net.Conn, rsaKey *rsa.PrivateKey, rsaPub *rsa.PublicKey) (*secureConn, error) {
	// Receive client's ephemeral DH public key
	b := make([]byte, pubKeySize)
	_, err := c.Read(b)
	if err != nil {
		return nil, err
	}

	remoteKey, err := ecdh.X25519().NewPublicKey(b)
	if err != nil {
		c.Close()
		return nil, err
	}

	// Generate our own ephemeral DH key pair
	privkey, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		c.Close()
		return nil, err
	}

	// Compute (s)hared (s)ecret
	ss, err := privkey.ECDH(remoteKey)
	if err != nil {
		c.Close()
		return nil, err
	}

	ss = genSharedKey(ss)

	pubkey := privkey.PublicKey()

	// Get byte representation of RSA public key to send
	rsaPubBytes := x509.MarshalPKCS1PublicKey(rsaPub)
	// Unsure if rsaPubLen is necessary, will research
	rsaPubLen := make([]byte, rsaKeyByteLen)
	binary.BigEndian.PutUint16(rsaPubLen, uint16(len(rsaPubBytes)))

	// Sign all data so far with our RSA private key
	opts := rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
		Hash:       crypto.BLAKE2b_512,
	}

	hash := opts.HashFunc().New()
	_, err = hash.Write(slices.Concat(remoteKey.Bytes(), pubkey.Bytes(), rsaPubBytes))
	if err != nil {
		c.Close()
		return nil, err
	}

	sig, err := rsaKey.Sign(rand.Reader, hash.Sum(nil), &opts)
	if err != nil {
		c.Close()
		return nil, err
	}

	sigLen := make([]byte, 2)
	binary.BigEndian.PutUint16(sigLen, uint16(len(sig)))

	// Generate a MAC over all data so far using the DH shared secret
	// I may be wrong, but I belive that since I'm using XChaCha20-Poly1305 it should be sufficient
	// to encrypt the data as normal to serve as a MAC
	mac, err := encryptBytes(slices.Concat(remoteKey.Bytes(), pubkey.Bytes(), rsaPubBytes, sig), ss, nil)
	if err != nil {
		c.Close()
		return nil, err
	}

	macLen := make([]byte, 2)
	binary.BigEndian.PutUint16(macLen, uint16(len(mac)))

	// Send server hello
	_, err = c.Write(slices.Concat(pubkey.Bytes(), rsaPubLen, rsaPubBytes, sigLen, sig, macLen, mac))
	if err != nil {
		c.Close()
		return nil, err
	}

	return &secureConn{c, ss, structures.Queue[byte]{}}, nil
}

// TODO: Chunking for oversize packets
func (s *secureConn) Read(b []byte) (int, error) {
	n := 0
	for s.queue.HasNext() && n < len(b) {
		b[n] = s.queue.Dequeue()
		n++
	}

	// Don't bother reading more from s.c if the supplied buffer filled up from the queue
	if n < len(b) {
		// Max size is too big, should be uint16 rather than uint32
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

		d, err := decryptBytes(buf, s.ss, nil)
		if err != nil {
			return 0, err
		}

		for i := 0; i < len(d); i++ {
			if n < len(b) {
				b[n] = d[i]
				n++
			} else {
				s.queue.Enqueue(d[i])
			}
		}
	}

	return n, nil
}

func (s *secureConn) Write(b []byte) (int, error) {
	e, err := encryptBytes(b, s.ss, nil)
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
