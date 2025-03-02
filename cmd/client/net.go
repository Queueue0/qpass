package main

import (
	"errors"
	"log"
	"net"
	"time"

	"github.com/Queueue0/qpass/internal/crypto"
	"github.com/Queueue0/qpass/internal/protocol"
	"github.com/google/uuid"
)

var (
	ErrNoActiveUser = errors.New("no logged in user")
	ErrPingFail     = errors.New("Unable to ping sync server")
	ErrCommFail     = errors.New("Communication with server failed unexpectedly")
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

func (app *Application) newUserSync(authToken []byte) (string, error) {
	nud := protocol.NewUserData{Token: authToken}
	b, err := nud.Encode()
	if err != nil {
		return "", err
	}

	c, err := crypto.Dial(app.ServerAddress)
	if err != nil {
		return "", ErrPingFail
	}
	defer c.Close()

	protocol.NewPing().WriteTo(c)
	r := protocol.Payload{}
	r.ReadFrom(c)

	if r.Type() != protocol.PONG {
		return "", ErrPingFail
	}

	// Will never error here, only possible error is if max payload size is
	// exceeded which won't happen with just an auth token
	p, _ := protocol.NewPayload(protocol.NUSR, b)
	p.WriteTo(c)

	r = protocol.Payload{}
	r.ReadFrom(c)

	if r.Type() == protocol.FAIL {
		return "", errors.New(string(r.Bytes()))
	}

	if r.Type() != protocol.SUCC {
		return "", ErrCommFail
	}

	id := string(r.Bytes())
	_, err = uuid.Parse(id)

	protocol.NewSucc().WriteTo(c)
	return id, err
}

func (app *Application) sync() error {
	if app.ActiveUser == nil {
		return ErrNoActiveUser
	}

	activeUUID := app.ActiveUser.ID.String()
	ad := protocol.AuthData{Token: app.ActiveUser.AuthToken}
	authBytes, err := ad.Encode()
	if err != nil {
		return err
	}

	apl, err := protocol.NewPayload(protocol.AUTH, authBytes)
	if err != nil {
		return err
	}

	c, err := crypto.Dial(app.ServerAddress)
	if err != nil {
		return err
	}
	defer c.Close()

	_, err = apl.WriteTo(c)
	if err != nil {
		return err
	}

	authResp := protocol.Payload{}
	_, err = authResp.ReadFrom(c)
	if err != nil {
		return err
	}

	if authResp.Type() == protocol.FAIL {
		protocol.NewSucc().WriteTo(c)
		return errors.New("Remote error: " + authResp.String())
	}

	if authResp.Type() != protocol.SUCC {
		return errors.New("Unexpected response type from server")
	}

	lastSync, err := app.Logs.GetLastSync(activeUUID)
	if err != nil {
		return err
	}

	logs, err := app.Logs.GetAllSince(lastSync, activeUUID)
	if err != nil {
		return err
	}

	sd := protocol.SyncData{LastSync: lastSync, UUID: activeUUID, Logs: logs}
	bytes, err := sd.Encode()
	if err != nil {
		return err
	}

	pl, err := protocol.NewPayload(protocol.SYNC, bytes)
	if err != nil {
		return err
	}

	_, err = pl.WriteTo(c)
	if err != nil {
		return err
	}

	r := protocol.Payload{}
	_, err = r.ReadFrom(c)
	if err != nil {
		return err
	}
	if r.Type() == protocol.FAIL {
		// For now, SUCC is the only way to terminate a connection
		// Should probably add a DONE or something instead
		protocol.NewSucc().WriteTo(c)
		return errors.New("Remote error: " + r.String())
	}

	rd := protocol.SyncData{}
	err = rd.Decode(r.Bytes())
	if err != nil {
		return err
	}

	// TODO: Handle response data
	for _, l := range rd.Logs {
		log.Println(l.String())
	}

	err = app.Logs.SetLastSync(time.Now(), activeUUID)
	if err != nil {
		log.Println(err.Error())
	}

	protocol.NewSucc().WriteTo(c)

	return nil
}

func (app *Application) syncUsers() error {
	if app.ActiveUser == nil {
		return ErrNoActiveUser
	}

	lastSync, err := app.Logs.GetLastSync(app.ActiveUser.ID.String())
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
		return errors.New("Remote error: " + response.String())
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
