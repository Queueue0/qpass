package protocol

import (
	"bytes"
	"encoding/gob"

	"github.com/Queueue0/qpass/internal/models"
)

type PayloadData interface {
	Encode() ([]byte, error)
	Decode([]byte)
}

type SyncData struct {
	UUID      string
	Passwords models.PasswordList
}

func (s *SyncData) Encode() (data []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err = enc.Encode(s)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *SyncData) Decode(data []byte) error {
	var buf bytes.Buffer
	_, err := buf.Write(data)
	if err != nil {
		return err
	}

	dec := gob.NewDecoder(&buf)
	err = dec.Decode(s)
	if err != nil {
		return err
	}

	return nil
}

type AuthData struct {
	Token []byte
}

func (d *AuthData) Encode() (data []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err = enc.Encode(d)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (d *AuthData) Decode(data []byte) error {
	var buf bytes.Buffer
	_, err := buf.Write(data)
	if err != nil {
		return err
	}

	dec := gob.NewDecoder(&buf)
	err = dec.Decode(d)
	if err != nil {
		return err
	}

	return nil
}

type NewUserData struct {
	UUID  string
	Token []byte
}

func (d *NewUserData) Encode() (data []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err = enc.Encode(d)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (d *NewUserData) Decode(data []byte) error {
	var buf bytes.Buffer
	_, err := buf.Write(data)
	if err != nil {
		return err
	}

	dec := gob.NewDecoder(&buf)
	err = dec.Decode(d)
	if err != nil {
		return err
	}

	return nil
}
