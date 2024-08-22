package models

import (
	"database/sql"

	"github.com/Queueue0/qpass/internal/crypto"
)

// e means encrypted
type Password struct {
	ID           int
	UserID       int
	ServiceName  string
	Username     string
	Password     string
	eServiceName string
	eUsername    string
	ePassword    string
}

func (p *Password) decrypt(u User) error {
	var err error
	p.ServiceName, err = crypto.Decrypt(p.eServiceName, u.Key)
	if err != nil {
		return err
	}

	p.Username, err = crypto.Decrypt(p.eUsername, u.Key)
	if err != nil {
		return err
	}

	p.ePassword, err = crypto.Decrypt(p.ePassword, u.Key)
	if err != nil {
		return err
	}

	return nil
}

type PasswordModel struct {
	DB *sql.DB
}

func (m *PasswordModel) Insert(u User, serviceName, username, password string) (int, error) {
	eServiceName, err := crypto.Encrypt(serviceName, u.Key)
	if err != nil {
		return 0, nil
	}

	eUsername, err := crypto.Encrypt(username, u.Key)
	if err != nil {
		return 0, nil
	}

	ePassword, err := crypto.Encrypt(password, u.Key)
	if err != nil {
		return 0, nil
	}

	stmt := `INSERT INTO passwords (userID, service, username, password) VALUES (?, ?, ?, ?)`
	result, err := m.DB.Exec(stmt, u.ID, eServiceName, eUsername, ePassword)
	if err != nil {
		return 0, nil
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return int(id), nil
}

func (m *PasswordModel) GetAllForUser(u User) ([]Password, error) {
	return []Password{}, nil
}

func (m *PasswordModel) Search(u User, searchTerm string) ([]Password, error) {
	return []Password{}, nil
}
