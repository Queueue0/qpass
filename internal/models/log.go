package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type LogType string

// Possible Change Types
const (
	// Add user
	AUSR LogType = "AUSR"
	// Modify user
	MUSR LogType = "MUSR"
	// Delete user
	DUSR LogType = "DUSR"
	// Add password
	APWD LogType = "APWD"
	// Modify password
	MPWD LogType = "MPWD"
	// Delete password
	DPWD LogType = "DPWD"
)

var (
	LogWriteError = errors.New("Error writing to log")
)

type LogModel struct {
	DB *sql.DB
}

type Log struct {
	Timestamp  time.Time
	Type       LogType
	User       string
	OldService string
	NewService string
	OldName    string
	NewName    string
	OldPW      string
	NewPW      string
}

func (l Log) String() string {
	return fmt.Sprintf("%s %s %s %s %s %s %s %s %s", l.Timestamp.String(), l.Type, l.User, l.OldService, l.NewService, l.OldName, l.NewName, l.OldPW, l.NewPW)
}

func (l Log) Write(db *sql.DB) error {
	stmt := `INSERT INTO log (
		change_type, user,
		old_service, new_service,
		old_name, new_name,
		old_password, new_password
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := db.Exec(stmt, l.Type, l.User, l.OldService, l.NewService, l.OldName, l.NewName, l.OldPW, l.NewPW)
	if err != nil {
		return err
	}

	return nil
}

func (m *LogModel) GetAll() ([]Log, error) {
	stmt := `SELECT timestamp, change_type, user, old_service, new_service, old_name, new_name, old_password, new_password FROM log`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	ls := []Log{}
	for rows.Next() {
		l := Log{}
		err = rows.Scan(&l.Timestamp, &l.Type, &l.User, &l.OldService, &l.NewService, &l.OldName, &l.NewName, &l.OldPW, &l.NewPW)
		if err != nil {
			return nil, err
		}

		ls = append(ls, l)
	}

	return ls, nil
}

func (m *LogModel) GetAllSince(t time.Time, user string) ([]Log, error) {
	stmt := `SELECT timestamp, change_type, user, old_service, new_service, old_name, new_name, old_password, new_password FROM log
	WHERE timestamp > ? AND user = ?`

	rows, err := m.DB.Query(stmt, t, user)
	if err != nil {
		return nil, err
	}

	ls := []Log{}
	for rows.Next() {
		l := Log{}
		err = rows.Scan(&l.Timestamp, &l.Type, &l.User, &l.OldService, &l.NewService, &l.OldName, &l.NewName, &l.OldPW, &l.NewPW)
		if err != nil {
			return nil, err
		}

		ls = append(ls, l)
	}

	return ls, nil
}

func (m *LogModel) GetAllUserSince(t time.Time) ([]Log, error) {
	stmt := `SELECT timestamp, change_type, user, old_service, new_service, old_name, new_name, old_password, new_password FROM log
	WHERE timestamp > ? AND change_type LIKE "_USR"`

	rows, err := m.DB.Query(stmt, t)
	if err != nil {
		return nil, err
	}

	ls := []Log{}
	for rows.Next() {
		l := Log{}
		err = rows.Scan(&l.Timestamp, &l.Type, &l.User, &l.OldService, &l.NewService, &l.OldName, &l.NewName, &l.OldPW, &l.NewPW)
		if err != nil {
			return nil, err
		}

		ls = append(ls, l)
	}

	return ls, nil
}

func (m *LogModel) GetAllPasswordSince(t time.Time, u string) ([]Log, error) {
	stmt := `SELECT timestamp, change_type, user, old_service, new_service, old_name, new_name, old_password, new_password FROM log
	WHERE timestamp > ? AND change_type LIKE "_PWD" AND user = ?`

	rows, err := m.DB.Query(stmt, t, u)
	if err != nil {
		return nil, err
	}

	ls := []Log{}
	for rows.Next() {
		l := Log{}
		err = rows.Scan(&l.Timestamp, &l.Type, &l.User, &l.OldService, &l.NewService, &l.OldName, &l.NewName, &l.OldPW, &l.NewPW)
		if err != nil {
			return nil, err
		}

		ls = append(ls, l)
	}

	return ls, nil
}

func (m *LogModel) GetLastSync(uid string) (time.Time, error) {
	stmt := `SELECT timestamp FROM last_user_sync WHERE user=?`

	row := m.DB.QueryRow(stmt, uid)
	var t time.Time
	err := row.Scan(&t)

	return t, err
}

func (m *LogModel) SetLastSync(t time.Time, uid string) error {
	stmt := `UPDATE last_user_sync SET timestamp=? WHERE user=?`

	_, err := m.DB.Exec(stmt, t, uid)
	return err
}

func (m *LogModel) NewLastSync(uid string) error {
	t := time.Time{}
	stmt := `INSERT INTO last_user_sync (user, timestamp) VALUES (?, ?)`

	_, err := m.DB.Exec(stmt, uid, t)
	return err
}
