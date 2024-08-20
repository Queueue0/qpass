package models

import "database/sql"

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

type PasswordModel struct {
	DB *sql.DB
}

func (m *PasswordModel) Insert(serviceName, username, password string) (int, error) {
	return 0, nil
}

func (m *PasswordModel) GetAllForUser(userId int) ([]Password, error) {
	return []Password{}, nil
}

func (m *PasswordModel) Search(userId int, searchTerm string) ([]Password, error) {
	return []Password{}, nil
}
