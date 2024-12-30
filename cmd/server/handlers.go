package main

import (
	"log"
	"net"

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

func (app *Application) userSync(p protocol.Payload, c net.Conn) {
	var sd protocol.SyncData
	err := sd.Decode(p.Bytes())
	if err != nil {
		log.Println(err.Error())
		fail := protocol.NewFail(err.Error())
		fail.WriteTo(c)
		return
	}

	for _, l := range sd.Logs {
		log.Println(l.String())
	}

	// For now, just respond with the same payload
	// TODO: Make this actually work
	p.WriteTo(c)
}
