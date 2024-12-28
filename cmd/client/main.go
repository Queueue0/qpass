package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"gioui.org/app"
	"gioui.org/unit"
	"github.com/Queueue0/qpass/internal/dbman"
	"github.com/Queueue0/qpass/internal/models"
	"github.com/Queueue0/qpass/internal/protocol"
)

type Application struct {
	UserModel     *models.UserModel
	ActiveUser    *models.User
	PasswordModel *models.PasswordModel
	Passwords     models.PasswordList
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

	c, err := net.Dial("tcp", "127.0.0.1:10448")
	if err != nil {
		log.Println("Failed to connect to server")
	}

	if c != nil {
		defer c.Close()
		ping := protocol.NewPing()
		ping.WriteTo(c)
		var response protocol.Payload
		response.ReadFrom(c)

		if response.Type() == protocol.PONG {
			log.Println("PONG")
		} else {
			log.Println("Ping failed")
		}
	}

	go func() {
		if a.UserModel.Count() <= 0 {
			created := false
			aw := new(app.Window)
			aw.Option(app.Title("New User"))
			aw.Option(app.Size(unit.Dp(1280), unit.Dp(720)))
			created, err = a.newUserView(aw)
			if err != nil {
				fmt.Println(err.Error())
			}

			if !created {
				os.Exit(0)
			}
		}

		lw := new(app.Window)
		lw.Option(app.Title("Login"))
		lw.Option(app.Size(unit.Dp(500), unit.Dp(200)))
		if err := a.loginView(lw); err != nil {
			log.Fatal(err)
		}

		if a.ActiveUser == nil {
			os.Exit(0)
		}

		w := new(app.Window)
		w.Option(app.Title("QPass"))
		w.Option(app.Size(unit.Dp(1280), unit.Dp(720)))
		if err := a.mainView(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}
