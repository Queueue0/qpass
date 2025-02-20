package models

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/Queueue0/qpass/internal/crypto"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	encryptedUsername string
	salt              string
	ID                uuid.UUID
	Username          string
	Key               []byte
	AuthToken         []byte
}

func (u User) EUsername() string {
	return u.encryptedUsername
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

	uuid, err := uuid.NewRandom()
	if err != nil {
		return 0, err
	}

	stmt := `INSERT INTO users (uuid, username, salt) VALUES (?, ?, ?)`
	result, err := m.DB.Exec(stmt, uuid.String(), encryptedUsername, salt)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	l := Log{
		Type:    AUSR,
		User:    uuid.String(),
		OldName: "",
		NewName: encryptedUsername,
		OldPW:   "",
		NewPW:   salt,
	}

	err = l.Write(m.DB)
	if err != nil {
		fmt.Println(err.Error())
		return int(id), LogWriteError
	}

	return int(id), nil
}

func (m *UserModel) ServerInsert(u User) (int, error) {
	token := base64.RawStdEncoding.EncodeToString(u.AuthToken)
	stmt := `INSERT INTO users (uuid, auth_token) VALUES (?, ?)`
	result, err := m.DB.Exec(stmt, u.ID.String(), token)
	if err != nil {
		return 0, err
	}
	
	id, err := result.LastInsertId()
	return int(id), err
}

func (m *UserModel) GetByUUID(id string) (*User, error) {
	row := m.DB.QueryRow("SELECT uuid, auth_token FROM users WHERE uuid = ?", id)

	var uuidStr, tokenStr string
	err := row.Scan(&uuidStr, &tokenStr)
	if err != nil {
		return nil, err
	}

	var u User
	u.ID, err = uuid.Parse(uuidStr)
	if err != nil {
		return nil, err
	}

	u.AuthToken, err = base64.RawStdEncoding.DecodeString(tokenStr)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (m *UserModel) InsertFromLog(l Log) (int, error) {
	stmt := `INSERT INTO users (uuid, username, salt) VALUES (?, ?, ?)`
	result, err := m.DB.Exec(stmt, l.User, l.NewName, l.NewPW)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return int(id), err
}

func (m *UserModel) IDtoUUID(id int) (string, error) {
	stmt := `SELECT uuid FROM users WHERE id=?`
	row := m.DB.QueryRow(stmt, id)
	var UUID string
	err := row.Scan(&UUID)
	return UUID, err
}

func (m *UserModel) Authenticate(username, password string) (User, error) {
	stmt := `SELECT uuid, username, salt FROM users`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return User{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var u User
		var uuidString string

		err := rows.Scan(&uuidString, &u.encryptedUsername, &u.salt)
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
			u.ID, err = uuid.Parse(uuidString)
			if err != nil {
				return User{}, err
			}

			u.AuthToken = crypto.Hash([]byte(password), []byte(u.Username), 10)
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
