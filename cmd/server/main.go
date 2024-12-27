package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Queueue0/qpass/internal/dbman"
	"github.com/Queueue0/qpass/internal/models"
	"github.com/Queueue0/qpass/internal/protocol"
)

type Application struct {
	UserModel     *models.UserModel
	PasswordModel *models.PasswordModel
}

func main() {
	qpassHome, err := dbman.GetQpassHome()
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("file:%s/pwdb.sqlite?mode=rwc", qpassHome)
	db, err := dbman.OpenDB(dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = dbman.InitializeDB(db)
	if err != nil {
		log.Fatal(err)
	}

	um := models.UserModel{
		DB: db,
	}

	pm := models.PasswordModel{
		DB: db,
	}

	a := Application{
		UserModel:     &um,
		PasswordModel: &pm,
	}

	srv, err := net.Listen("tcp", "127.0.0.1:10448")
	if err != nil {
		panic(err)
	}
	defer srv.Close()

	for {
		c, err := srv.Accept()
		if err != nil {
			panic(err)
		}

		go a.respond(c)
	}
}

func (app *Application) respond(c net.Conn) {
	defer c.Close()
	p := protocol.Read(c)

	switch p.Type() {
	case protocol.PING:
		protocol.Write(c, protocol.NewPong())
	case protocol.SYNC:
		app.sync(p)
	}
}
