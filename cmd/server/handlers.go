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

	for _, p := range sd.Passwords {
		exists, err := app.passwords.Exists(p.UUID.String())
		if err != nil {
			protocol.NewFail(err.Error()).WriteTo(c)
			return
		}

		if !exists {
			err = app.passwords.DumbInsert(p)
			if err != nil {
				protocol.NewFail(err.Error()).WriteTo(c)
				return
			}
			continue
		}

		if p.Deleted {
			err = app.passwords.Delete(p.UUID.String())
			if err != nil {
				protocol.NewFail(err.Error()).WriteTo(c)
				return
			}
			continue
		}

		current, err := app.passwords.GetByUUID(p.UUID.String())
		if err != nil {
			protocol.NewFail(err.Error()).WriteTo(c)
			return
		}

		if p.LastChanged.After(current.LastChanged) {
			err = app.passwords.DumbUpdate(p)
			if err != nil {
				protocol.NewFail(err.Error()).WriteTo(c)
				return
			}
		}
	}

	id, err := uuid.Parse(sd.UUID)
	if err != nil {
		protocol.NewFail(err.Error()).WriteTo(c)
		return
	}

	pws, err := app.passwords.GetAllEncryptedForUser(models.User{ID: id})
	if err != nil {
		protocol.NewFail(err.Error()).WriteTo(c)
		return
	}

	rd := protocol.SyncData{
		Passwords: pws,
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
}

func (app *Application) authenticate(p protocol.Payload) (bool, string, error) {
	var ad protocol.AuthData
	err := ad.Decode(p.Bytes())
	if err != nil {
		return false, "", err
	}

	ad.Token = crypto.ServerAuthToken(ad.Token)

	u, err := app.users.ServerGetByAuthToken(ad.Token)
	if err != nil {
		log.Println(err.Error())
		return false, "", err
	}

	return bytes.Equal(ad.Token, u.AuthToken), u.ID.String(), nil
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

	nud.Token = crypto.ServerAuthToken(nud.Token)

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
		temp, err := uuid.NewRandom()
		if err != nil {
			protocol.NewFail(ErrUserCreateFail.Error())
			return err
		}
		nud.UUID = temp.String()
	}

	// Will never panic because of the above validation
	UUID := uuid.MustParse(nud.UUID)

	_, err = app.users.ServerInsert(models.User{ID: UUID, AuthToken: nud.Token})
	if err != nil {
		protocol.NewFail(ErrUserCreateFail.Error()).WriteTo(c)
		return err
	}

	_, err = protocol.NewSuccWithData([]byte(nud.UUID)).WriteTo(c)
	return err
}
