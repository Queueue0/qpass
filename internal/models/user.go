package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	encryptedUsername string
	salt              string
	ID                string
	Username          string
}

type UserModel struct {
	DB sql.DB
}


