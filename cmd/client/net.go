package main

import (
	"log"
	"net"

	"github.com/Queueue0/qpass/internal/crypto"
	"github.com/Queueue0/qpass/internal/protocol"
)

func (app *Application) send(p *protocol.Payload) error {
	c, err := net.Dial("tcp", app.ServerAddress)
	if err != nil {
		return err
	}

	defer c.Close()
	_, err = p.WriteTo(c)

	return err
}

func (app *Application) sync() error {
	return nil
}

func (app *Application) syncUsers() error {
	lastSync, err := app.Logs.GetLastSync()
	if err != nil {
		return err
	}

	logs, err := app.Logs.GetAllUserSince(lastSync)
	if err != nil {
		return err
	}

	sd := protocol.SyncData{LastSync: lastSync, Logs: logs}
	data, err := sd.Encode()
	if err != nil {
		return err
	}

	payload, err := protocol.NewPayload(protocol.SUSR, data)
	if err != nil {
		return err
	}

	c, err := crypto.Dial(app.ServerAddress)
	if err != nil {
		return err
	}

	defer c.Close()
	_, err = payload.WriteTo(c)

	if err != nil {
		return err
	}

	// TODO: Properly handle response
	response := protocol.Payload{}
	response.ReadFrom(c)

	if response.Type() == protocol.FAIL {
		log.Println(response.String())
		protocol.NewSucc().WriteTo(c)
	}

	responseData := protocol.SyncData{}
	err = responseData.Decode(response.Bytes())
	if err != nil {
		return err
	}

	for _, l := range responseData.Logs {
		log.Println(l.String())
	}

	protocol.NewSucc().WriteTo(c)

	return nil
}
