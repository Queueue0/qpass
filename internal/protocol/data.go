package protocol

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/Queueue0/qpass/internal/models"
)

type PayloadData interface {
	Encode() ([]byte, error)
	Decode([]byte)
}

type SyncData struct {
	LastSync time.Time
	UUID     string
	Logs     []models.Log
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
