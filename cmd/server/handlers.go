package main

import (
	"log"

	"github.com/Queueue0/qpass/internal/protocol"
)

func (app *Application) sync(p protocol.Payload) {
	var sd protocol.SyncData
	err := sd.Decode(p.Bytes())
	if err != nil {
		log.Println(err.Error())
		return
	}

	for _, l := range sd.Logs {
		log.Println(l.String())
	}
	return
}
