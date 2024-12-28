package models

import (
	"database/sql"
	"errors"

	"github.com/Queueue0/qpass/internal/crypto"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	encryptedUsername string
	salt              string
	ID                int
	Username          string
	Key               []byte
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(username, password string) (int, error) {
	salt, err := crypto.GenSalt(16)
	if err != nil {
		return 0, err
	}

	key := crypto.GetKey(password, salt)

	encryptedUsername, err := crypto.Encrypt(username, key)
	if err != nil {
		return 0, err
	}

	stmt := `INSERT INTO users (username, salt) VALUES (?, ?)`
	result, err := m.DB.Exec(stmt, encryptedUsername, salt)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	l := Log{
		Type:    AUSR,
		User:    encryptedUsername,
		OldName: "",
		NewName: encryptedUsername,
		OldPW:   "",
		NewPW:   salt,
	}

	err = l.Write(m.DB)
	if err != nil {
		return int(id), LogWriteError
	}

	return int(id), nil
}

func (m *UserModel) Authenticate(username, password string) (User, error) {
	stmt := `SELECT id, username, salt FROM users`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return User{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var u User

		err := rows.Scan(&u.ID, &u.encryptedUsername, &u.salt)
		if err != nil {
			return User{}, err
		}

		key := crypto.GetKey(password, u.salt)
		u.Username, err = crypto.Decrypt(u.encryptedUsername, key)
		if err != nil && err.Error() != "chacha20poly1305: message authentication failed" {
			return User{}, err
		}

		if u.Username == username {
			u.Key = key
			return u, nil
		}
	}

	return User{}, errors.New("Username or password is incorrect")
}

func (m *UserModel) Count() int {
	row := m.DB.QueryRow("SELECT COUNT(id) FROM users")
	var c int
	row.Scan(&c)
	return c
}
