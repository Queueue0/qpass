package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Queueue0/qpass/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	username, passwd := os.Args[1], os.Args[2]

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

	um := models.UserModel{
		DB: db,
	}

	pm := models.PasswordModel{
		DB: db,
	}

	u, err := um.Authenticate(username, passwd)
	if err != nil {
		log.Fatal(err)
	}

	_, err = pm.Insert(u, "testPlatform", "testUser", "testPassword")
	if err != nil {
		log.Fatal(err)
	}

	_, err = pm.Insert(u, "google", "testuser@gmail.com", "testgmail")
	if err != nil {
		log.Fatal(err)
	}
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
