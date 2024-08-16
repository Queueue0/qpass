package crypto

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
	chacha "golang.org/x/crypto/chacha20poly1305"
)

func Encrypt(key []byte, s string) (string, error) {
	aead, err := chacha.NewX(key)
	if err != nil {
		return "", err
	}

	sbytes := []byte(s)
	nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(sbytes)+aead.Overhead())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	encrypted := aead.Seal(nonce, nonce, sbytes, nil)

	encStr := base64.RawStdEncoding.EncodeToString(encrypted)

	return encStr, nil
}

func Decrypt(key []byte, s string) (string, error) {
	b, err := base64.RawStdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	aead, err := chacha.NewX(key)
	if err != nil {
		return "", err
	}

	nonce, enc := b[:aead.NonceSize()], b[aead.NonceSize():]

	text, err := aead.Open(nil, nonce, enc, nil)
	if err != nil {
		return "", err
	}

	return string(text), nil
}

func GetKey(password, salt string) []byte {
	key := argon2.IDKey([]byte(password), []byte(salt), 3, 64 * 1024, 2, 32)
	return key
}

func GenSalt(n uint32) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(b), nil
}
