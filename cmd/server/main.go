package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Queueue0/qpass/internal/crypto"
	"github.com/Queueue0/qpass/internal/dbman"
	"github.com/Queueue0/qpass/internal/models"
	"github.com/Queueue0/qpass/internal/protocol"
)

type Application struct {
	users     *models.UserModel
	passwords *models.PasswordModel
	homeDir   string
}

func main() {
	qpassHome, err := dbman.GetQpassServerHome()
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("file:%s/pwdb.sqlite?mode=rwc", qpassHome)
	db, err := dbman.OpenDB(dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = dbman.InitializeDB(db, false)
	if err != nil {
		log.Fatal(err)
	}

	if !haveKeys(qpassHome) {
		err = genKeyPair(qpassHome)
		if err != nil {
			log.Fatal(err)
		}
	}

	um := models.UserModel{
		DB: db,
	}

	pm := models.PasswordModel{
		DB: db,
	}

	a := Application{
		users:     &um,
		passwords: &pm,
		homeDir:   qpassHome,
	}

	srv, err := net.Listen("tcp", "127.0.0.1:10448")
	if err != nil {
		panic(err)
	}
	defer srv.Close()

	log.Println(fmt.Sprintf("Server started on %s", srv.Addr().String()))

	for {
		c, err := srv.Accept()
		if err != nil {
			panic(err)
		}

		go a.handle(c)
	}
}

func (app *Application) handle(c net.Conn) {
	log.Println("Received connection", c.RemoteAddr().String())
	kp, err := getKeyPair(app.homeDir)
	if err != nil {
		log.Println(err.Error())
		return
	}

	sc, err := crypto.NewServerConn(c, kp.key, kp.pubKey)
	if err != nil {
		log.Println(c.RemoteAddr(), err.Error())
		return
	}

	app.respond(sc)
}

func (app *Application) respond(c net.Conn) {
	defer c.Close()

connLoop:
	for {
		var p protocol.Payload
		_, err := p.ReadFrom(c)
		if err != nil {
			log.Println(c.RemoteAddr(), err.Error())
			return
		}

		log.Println(c.RemoteAddr(), p.TypeString())

		switch p.Type() {
		case protocol.PING:
			protocol.NewPong().WriteTo(c)
		case protocol.SYNC:
			app.sync(p)
		case protocol.SUSR:
			app.userSync(p, c)
		case protocol.SUCC:
			break connLoop
		}
	}
}
