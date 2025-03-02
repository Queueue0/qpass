package main

import (
	"bytes"
	"errors"
	"log"
	"net"

	"github.com/Queueue0/qpass/internal/crypto"
	"github.com/Queueue0/qpass/internal/models"
	"github.com/Queueue0/qpass/internal/protocol"
	"github.com/google/uuid"
)

func (app *Application) sync(p protocol.Payload, c net.Conn) {
	var sd protocol.SyncData
	err := sd.Decode(p.Bytes())
	if err != nil {
		log.Println(err.Error())
		protocol.NewFail(err.Error()).WriteTo(c)
		return
	}

	responseLogs, err := app.logs.GetAllSince(sd.LastSync, sd.UUID)
	if err != nil {
		protocol.NewFail(err.Error()).WriteTo(c)
		return
	}

	rd := protocol.SyncData{
		Logs: responseLogs,
	}
	rdBytes, err := rd.Encode()
	if err != nil {
		protocol.NewFail(err.Error()).WriteTo(c)
		return
	}

	response, err := protocol.NewPayload(protocol.SYNC, rdBytes)
	if err != nil {
		protocol.NewFail(err.Error()).WriteTo(c)
		return
	}

	response.WriteTo(c)

	// write received logs
	for _, l := range sd.Logs {
		l.Write(app.users.DB)
	}
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

func (app *Application) authenticate(p protocol.Payload) (bool, error) {
	var ad protocol.AuthData
	err := ad.Decode(p.Bytes())
	if err != nil {
		return false, err
	}

	ad.Token = crypto.Hash(ad.Token, nil, 30)

	u, err := app.users.ServerGetByAuthToken(ad.Token)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}

	return bytes.Equal(ad.Token, u.AuthToken), nil
}

var (
	ErrUserExists     = errors.New("User already exists")
	ErrUserCreateFail = errors.New("Failed to create new user")
)

func (app *Application) newUser(p protocol.Payload, c net.Conn) error {
	var nud protocol.NewUserData
	err := nud.Decode(p.Bytes())
	if err != nil {
		protocol.NewFail(ErrUserCreateFail.Error()).WriteTo(c)
		return err
	}

	// Check if user with same auth token or UUID exists
	// if so, fail
	_, err = app.users.ServerGetByAuthToken(nud.Token)
	if err == nil {
		protocol.NewFail(ErrUserExists.Error()).WriteTo(c)
		return ErrUserExists
	}

	_, err = app.users.GetByUUID(nud.UUID)
	if err == nil {
		protocol.NewFail(ErrUserExists.Error()).WriteTo(c)
		return ErrUserExists
	}

	if err = uuid.Validate(nud.UUID); err != nil {
		// Might panic because Google's engineers are bastards I guess
		nud.UUID = uuid.New().String()
	}

	// Will never panic because of prior validation
	UUID := uuid.MustParse(nud.UUID)
	_, err = app.users.ServerInsert(models.User{ID: UUID, AuthToken: nud.Token})
	if err != nil {
		protocol.NewFail(ErrUserCreateFail.Error()).WriteTo(c)
		return err
	}

	_, err = protocol.NewSuccWithData([]byte(nud.UUID)).WriteTo(c)
	return err
}
