package crypto

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
	chacha "golang.org/x/crypto/chacha20poly1305"
)

func Encrypt(s string, key []byte) (string, error) {
	sbytes := []byte(s)
	encrypted, err := encryptBytes(sbytes, key, nil)
	if err != nil {
		return "", err
	}

	encStr := base64.RawStdEncoding.EncodeToString(encrypted)

	return encStr, nil
}

func Decrypt(s string, key []byte) (string, error) {
	b, err := base64.RawStdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	unsealed, err := decryptBytes(b, key, nil)
	if err != nil {
		return "", err
	}

	return string(unsealed), nil
}

func encryptBytes(plaintext, key, additionalData []byte) ([]byte, error) {
	aead, err := chacha.NewX(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(plaintext)+aead.Overhead())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	sealed := aead.Seal(nonce, nonce, plaintext, additionalData)

	return sealed, nil
}

func decryptBytes(cyphertext, key, additionalData []byte) ([]byte, error) {
	aead, err := chacha.NewX(key)
	if err != nil {
		return nil, err
	}

	nonce, enc := cyphertext[:aead.NonceSize()], cyphertext[aead.NonceSize():]

	unsealed, err := aead.Open(nil, nonce, enc, additionalData)
	if err != nil {
		return nil, err
	}

	return unsealed, nil
}

func GetKey(password, salt string) []byte {
	key := argon2.IDKey([]byte(password), []byte(salt), 3, 64*1024, 2, 32)
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
