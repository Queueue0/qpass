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

func (m *UserModel) Insert(username, password, uuid string) (int, error) {
	key := crypto.GetKey(password, username)

	encryptedUsername, err := crypto.Encrypt(username, key)
	if err != nil {
		return 0, err
	}

	stmt := `INSERT INTO users (uuid, username) VALUES (?, ?)`
	result, err := m.DB.Exec(stmt, uuid, encryptedUsername)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	l := Log{
		Type:    AUSR,
		User:    uuid,
		OldName: "",
		NewName: encryptedUsername,
		OldPW:   "",
		NewPW:   "",
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

	return parseFromStrings(uuidStr, tokenStr)
}

func (m *UserModel) ServerGetByAuthToken(at []byte) (*User, error) {
	t := base64.RawStdEncoding.EncodeToString(at)
	row := m.DB.QueryRow("SELECT uuid, auth_token FROM users WHERE auth_token = ?", t)

	var uuidStr, tokenStr string
	err := row.Scan(&uuidStr, &tokenStr)
	if err != nil {
		return nil, err
	}

	return parseFromStrings(uuidStr, tokenStr)
}

func parseFromStrings(uuidStr, tokenStr string) (*User, error) {
	var u User
	var err error
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
	stmt := `INSERT INTO users (uuid, username) VALUES (?, ?)`
	result, err := m.DB.Exec(stmt, l.User, l.NewName)
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

var ErrIncorrectCreds = errors.New("Username or password is incorrect")

func (m *UserModel) Authenticate(username, password string) (User, error) {
	stmt := `SELECT uuid, username FROM users`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return User{}, err
	}

	defer rows.Close()

	key := crypto.GetKey(password, username)
	for rows.Next() {
		var u User
		var uuidString string

		err := rows.Scan(&uuidString, &u.encryptedUsername)
		if err != nil {
			return User{}, err
		}

		u.Username, err = crypto.Decrypt(u.encryptedUsername, key)
		if err != nil {
			if err.Error() == "chacha20poly1305: message authentication failed" {
				return User{}, ErrIncorrectCreds
			}
			return User{}, err
		}

		if u.Username == username {
			u.Key = key
			u.ID, err = uuid.Parse(uuidString)
			if err != nil {
				return User{}, err
			}

			u.AuthToken = crypto.Hash(key, []byte(password), 10)
			return u, nil
		}
	}

	return User{}, ErrIncorrectCreds
}

func (m *UserModel) Count() int {
	row := m.DB.QueryRow("SELECT COUNT(id) FROM users")
	var c int
	row.Scan(&c)
	return c
}
