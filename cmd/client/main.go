package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"gioui.org/app"
	"gioui.org/unit"
	"github.com/Queueue0/qpass/internal/models"
	"github.com/Queueue0/qpass/internal/protocol"
	_ "github.com/mattn/go-sqlite3"
)

type Application struct {
	UserModel     *models.UserModel
	ActiveUser    *models.User
	PasswordModel *models.PasswordModel
	Passwords     models.PasswordList
}

func main() {
	qpassHome, err := getQpassHome()
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("file:%s/pwdb.sqlite?mode=rwc", qpassHome)
	db, err := openDB(dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = initializeDB(db)
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
		UserModel: &um,
		PasswordModel: &pm,
	}

	c, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Println("Failed to connect to server")
	}
	defer c.Close()

	protocol.Write(c, protocol.NewPing())
	p := protocol.Read(c)

	if p.Type() == protocol.PONG {
		log.Println("PONG")
	} else {
		log.Println("Ping failed")
	}

	go func() {
		lw := new(app.Window)
		lw.Option(app.Title("Login"))
		lw.Option(app.Size(unit.Dp(500), unit.Dp(200)))
		if err := a.loginView(lw); err != nil {
			log.Fatal(err)
		}

		if a.ActiveUser != nil {
			w := new(app.Window)
			w.Option(app.Title("QPass"))
			w.Option(app.Size(unit.Dp(1280), unit.Dp(720)))
			if err := a.mainView(w); err != nil {
				log.Fatal(err)
			}
		}
		os.Exit(0)
	}()

	app.Main()
}

func getQpassHome() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	qpassHome := fmt.Sprintf("%s/.qpass", homeDir)
	if _, err := os.Stat(qpassHome); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(qpassHome, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	return qpassHome, nil
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func initializeDB(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT, salt TEXT)")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS passwords (id INTEGER PRIMARY KEY, userId INT, service TEXT, username TEXT, password TEXT)")
	if err != nil {
		return err
	}

	return nil
}
