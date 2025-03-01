package dbman

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

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
	if client {
		stmt = "CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, uuid TEXT UNIQUE, username TEXT, last_sync DATETIME)"
	} else {
		stmt = "CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, uuid TEXT UNIQUE, auth_token TEXT)"
	}
	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}

	// Server should only have to store logs, not passwords
	// Maybe should change to simplify syncing?
	if client {
		_, err = db.Exec("CREATE TABLE IF NOT EXISTS passwords (id INTEGER PRIMARY KEY, userId TEXT, service TEXT, username TEXT, password TEXT)")
		if err != nil {
			return err
		}
	}

	logstmt := `CREATE TABLE IF NOT EXISTS log (
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		change_type CHAR(4),
		user TEXT,
		old_service TEXT,
		new_service TEXT,
		old_name TEXT,
		new_name TEXT,
		old_password TEXT,
		new_password TEXT
	)`
	_, err = db.Exec(logstmt)
	if err != nil {
		return err
	}

	if client {
		_, err = db.Exec("CREATE TABLE IF NOT EXISTS last_user_sync (user TEXT, timestamp DATETIME)")
		if err != nil {
			return err
		}

		var c int
		r := db.QueryRow("SELECT COUNT(timestamp) FROM last_user_sync")
		err = r.Scan(&c)
		if err != nil {
			return err
		}

		if c <= 0 {
			t := time.Time{}
			_, err = db.Exec("INSERT INTO last_user_sync (user, timestamp) SELECT uuid, ? FROM users", t)
			if err != nil {
				return err
			}
		}
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
