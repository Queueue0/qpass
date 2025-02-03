package main

import (
	"fmt"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/unit"
	"github.com/Queueue0/qpass/internal/crypto"
	"github.com/Queueue0/qpass/internal/dbman"
	"github.com/Queueue0/qpass/internal/models"
	"github.com/Queueue0/qpass/internal/protocol"
)

type Application struct {
	UserModel     *models.UserModel
	ActiveUser    *models.User
	PasswordModel *models.PasswordModel
	Passwords     models.PasswordList
	Logs          *models.LogModel
	ServerAddress string
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

	err = dbman.InitializeDB(db, true)
	if err != nil {
		log.Fatal(err)
	}

	um := models.UserModel{
		DB: db,
	}

	pm := models.PasswordModel{
		DB: db,
	}

	lm := models.LogModel{
		DB: db,
	}

	a := Application{
		UserModel:     &um,
		PasswordModel: &pm,
		Logs:          &lm,
		ServerAddress: "127.0.0.1:10448",
	}

	sc, err := crypto.Dial(a.ServerAddress)
	if err != nil {
		log.Println("Failed to connect to server", err.Error())
	}

	if sc != nil {
		// Closing sc closes c
		defer sc.Close()
		ping := protocol.NewPing()
		_, err = ping.WriteTo(sc)
		if err != nil {
			log.Println("Write error", err.Error())
		}
		var response protocol.Payload
		_, err = response.ReadFrom(sc)
		if err != nil {
			log.Println("Read error", err.Error())
		}

		if response.Type() == protocol.PONG {
			log.Println("PONG")
		} else {
			log.Println("Ping failed")
		}
		succ := protocol.NewSucc()
		succ.WriteTo(sc)
	}

	a.syncUsers()

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
