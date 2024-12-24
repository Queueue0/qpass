package models

import (
	"database/sql"
	"errors"
)

type LogType string

// Possible Change Types
const (
	AUSR LogType = "AUSR" // Add user
	MUSR LogType = "MUSR" // Modify user
	DUSR LogType = "DUSR" // Delete user
	// etc, but for passwords
	APWD LogType = "APWD"
	MPWD LogType = "MPWD"
	DPWD LogType = "DPWD"
)

var (
	LogWriteError = errors.New("Error writing to log")
)

type model interface {
	GetDB() *sql.DB
}

type Log struct {
	DB      *sql.DB
	Type    LogType
	User    string
	OldName string
	NewName string
	OldPW   string
	NewPW   string
}

func (l *Log) Write() error {
	stmt := `INSERT INTO log (
		change_type, user,
		old_name, new_name,
		old_password, new_password
	) VALUES (?, ?, ?, ?, ?, ?)`

	_, err := l.DB.Exec(stmt, l.Type, l.User, l.OldName, l.NewName, l.OldPW, l.NewPW)
	if err != nil {
		return err
	}

	return nil
}
