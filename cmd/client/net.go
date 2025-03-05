package main

import (
	"errors"
	"fmt"
	"net"

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
	defer func() {
		protocol.NewSucc().WriteTo(c)
		c.Close()
	}()

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
		return errors.New("Remote error: " + authResp.String())
	}

	if authResp.Type() != protocol.SUCC {
		return errors.New("Unexpected response type from server")
	}

	pws, err := app.PasswordModel.GetAllForUser(*app.ActiveUser, true)
	if err != nil {
		return err
	}

	sd := protocol.SyncData{UUID: activeUUID, Passwords: pws}
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
		return errors.New("Remote error: " + r.String())
	}

	rd := protocol.SyncData{}
	err = rd.Decode(r.Bytes())
	if err != nil {
		return err
	}

	for _, p := range rd.Passwords {
		fmt.Println(p.UUID.String())
	}

	err = app.PasswordModel.ReplaceAllForUser(app.ActiveUser.ID.String(), rd.Passwords)
	if err != nil {
		return err
	}

	return nil
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

func (app *Application) loginSync(username, password string) error {
	// Try to authenticate first
	u, err := app.UserModel.Authenticate(username, password)
	if err == nil {
		// Should maybe make a call to sync in this block
		app.ActiveUser = &u
		return nil
	}

	ad := protocol.AuthData{Token: crypto.ClientAuthToken(username, password)}
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
	defer func() {
		protocol.NewSucc().WriteTo(c)
		c.Close()
	}()

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

	idStr := string(authResp.Bytes())
	_, err = app.UserModel.Insert(username, password, idStr)
	if err != nil {
		return err
	}

	u, err = app.UserModel.Authenticate(username, password)
	if err != nil {
		return err
	}
	app.ActiveUser = &u

	err = app.sync()
	if err != nil {
		return errors.New("sync: " + err.Error())
	}

	return nil
}
