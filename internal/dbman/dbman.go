package dbman

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func InitializeDB(db *sql.DB, client bool) error {
	var stmt string
	// Only users table needs to be different between client and server
	if client {
		stmt = "CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, uuid TEXT UNIQUE, username TEXT)"
	} else {
		stmt = "CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, uuid TEXT UNIQUE, auth_token TEXT)"
	}
	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS passwords (id INTEGER PRIMARY KEY, uuid TEXT UNIQUE, userId TEXT, service TEXT, username TEXT, password TEXT, last_changed DATETIME DEFAULT CURRENT_TIMESTAMP, deleted BOOLEAN DEFAULT FALSE)")
	if err != nil {
		return err
	}

	return nil
}

func getHome(dirname string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	qpassHome := fmt.Sprintf("%s/%s", homeDir, dirname)
	if _, err := os.Stat(qpassHome); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(qpassHome, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	return qpassHome, nil
}

func GetQpassHome() (string, error) {
	return getHome(".qpass")
}

func GetQpassServerHome() (string, error) {
	return getHome(".qpass_server")
}
