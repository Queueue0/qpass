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

func InitializeDB(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT, salt TEXT)")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS passwords (id INTEGER PRIMARY KEY, userId INT, service TEXT, username TEXT, password TEXT)")
	if err != nil {
		return err
	}

	logstmt := `CREATE TABLE IF NOT EXISTS log (
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		change_type CHAR(4),
		user TEXT,
		old_name TEXT,
		new_name TEXT,
		old_password TEXT,
		new_password TEXT
	)`
	_, err = db.Exec(logstmt)
	if err != nil {
		return err
	}

	return nil
}

func GetQpassHome() (string, error) {
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
