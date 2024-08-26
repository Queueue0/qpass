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

type PasswordList []Password

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

// Probably not the best way to write this...
func (p *Password) isDecrypted() bool {
	return p.ServiceName != "" && p.Username != "" && p.Password != ""
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

func (m *PasswordModel) GetAllForUser(u User) (PasswordList, error) {
	stmt := `SELECT id, userID, service, username, password FROM passwords WHERE userID = ?`
	rows, err := m.DB.Query(stmt, u.ID)
	if err != nil {
		return []Password{}, err
	}

	pws := []Password{}
	for rows.Next() {
		pw := Password{}
		err := rows.Scan(&pw.ID, &pw.UserID, &pw.eServiceName, &pw.eUsername, &pw.ePassword)
		if err != nil {
			return []Password{}, err
		}

		err = pw.decrypt(u)
		if err != nil {
			return []Password{}, err
		}

		pws = append(pws, pw)
	}

	return pws, nil
}

func (pl PasswordList) Search(u User, searchTerm string) (PasswordList, error) {

	return []Password{}, nil
}
