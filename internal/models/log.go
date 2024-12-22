package models

import "database/sql"

// Possible Change Types
const (
	AUSR = "AUSR" // Add user
	MUSR = "MUSR" // Modify user
	DUSR = "DUSR" // Delete user
	// etc, but for passwords
	APWD = "APWD"
	MPWD = "MPWD"
	DPWD = "DPWD"
)

type model interface {
	GetDB() *sql.DB
}
