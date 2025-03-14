package models

import (
	"bytes"
	"database/sql"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/Queueue0/qpass/internal/crypto"
	"github.com/google/uuid"
)

// e = encrypted
type Password struct {
	ID           int
	UUID         uuid.UUID
	UserID       uuid.UUID
	ServiceName  string
	Username     string
	Password     string
	LastChanged  time.Time
	Deleted      bool
	EServiceName string
	EUsername    string
	EPassword    string
}

type PasswordList []Password

func (p *Password) decrypt(u User) error {
	var err error
	p.ServiceName, err = crypto.Decrypt(p.EServiceName, u.Key)
	if err != nil {
		return err
	}

	p.Username, err = crypto.Decrypt(p.EUsername, u.Key)
	if err != nil {
		return err
	}

	p.Password, err = crypto.Decrypt(p.EPassword, u.Key)
	if err != nil {
		return err
	}

	return nil
}

// Probably not the best way to write this...
func (p *Password) isDecrypted() bool {
	return p.ServiceName != "" && p.Username != "" && p.Password != ""
}

// Notably not Equals(), because I'm only interested in if
// the UUID is the same
func (p *Password) IsSame(other Password) bool {
	// UUIDs are just [16]byte with receiver functions
	return bytes.Equal(p.UUID[:], other.UUID[:])
}

type PasswordModel struct {
	DB *sql.DB
}

func (m *PasswordModel) Insert(u User, serviceName, username, password string) (int, error) {
	eServiceName, err := crypto.Encrypt(serviceName, u.Key)
	if err != nil {
		return 0, err
	}

	eUsername, err := crypto.Encrypt(username, u.Key)
	if err != nil {
		return 0, err
	}

	ePassword, err := crypto.Encrypt(password, u.Key)
	if err != nil {
		return 0, err
	}

	UUID, err := uuid.NewRandom()
	if err != nil {
		return 0, err
	}

	stmt := `INSERT INTO passwords (uuid, userId, service, username, password) VALUES (?, ?, ?, ?, ?)`
	result, err := m.DB.Exec(stmt, UUID.String(), u.ID.String(), eServiceName, eUsername, ePassword)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *PasswordModel) Update(u User, id, serviceName, username, password string) error {
	eServiceName, err := crypto.Encrypt(serviceName, u.Key)
	if err != nil {
		return err
	}

	eUsername, err := crypto.Encrypt(username, u.Key)
	if err != nil {
		return err
	}

	ePassword, err := crypto.Encrypt(password, u.Key)
	if err != nil {
		return err
	}

	stmt := `UPDATE passwords SET service = ?, username = ?, password = ? WHERE uuid = ?`
	_, err = m.DB.Exec(stmt, eServiceName, eUsername, ePassword, id)
	return err
}

func (m *PasswordModel) Get(id int, u User) (Password, error) {
	stmt := `SELECT id, uuid, userId, service, username, password, last_changed, deleted FROM passwords WHERE id = ?`
	r := m.DB.QueryRow(stmt, id)

	p := Password{}
	var uuidStr string
	err := r.Scan(&p.ID, &uuidStr, &p.UserID, &p.EServiceName, &p.EUsername, &p.EPassword, &p.LastChanged, &p.Deleted)
	if err != nil {
		return Password{}, err
	}

	p.UUID, err = uuid.Parse(uuidStr)
	if err != nil {
		return Password{}, err
	}

	err = p.decrypt(u)
	if err != nil {
		return Password{}, err
	}

	return p, nil
}

func (m *PasswordModel) GetByUUID(UUID string) (Password, error) {
	stmt := `SELECT id, uuid, userId, service, username, password, last_changed, deleted FROM passwords WHERE uuid = ?`
	r := m.DB.QueryRow(stmt, UUID)

	p := Password{}
	var uuidStr string
	var useridStr string
	err := r.Scan(&p.ID, &uuidStr, &useridStr, &p.EServiceName, &p.EUsername, &p.EPassword, &p.LastChanged, &p.Deleted)
	if err != nil {
		return Password{}, err
	}

	p.UUID, err = uuid.Parse(uuidStr)
	if err != nil {
		return Password{}, err
	}

	p.UserID, err = uuid.Parse(useridStr)
	if err != nil {
		return Password{}, err
	}

	// No need to decrypt in cases where this function is used
	return p, nil
}

func (m *PasswordModel) GetAllForUser(u User, includeDeleted bool) (PasswordList, error) {
	stmt := `SELECT id, uuid, userId, service, username, password, last_changed, deleted FROM passwords WHERE userId = ?`
	rows, err := m.DB.Query(stmt, u.ID.String())
	if err != nil {
		return nil, err
	}

	pws := PasswordList{}
	for rows.Next() {
		pw := Password{}
		var uuidString string
		var useridString string
		err := rows.Scan(&pw.ID, &uuidString, &useridString, &pw.EServiceName, &pw.EUsername, &pw.EPassword, &pw.LastChanged, &pw.Deleted)
		if err != nil {
			return nil, err
		}

		if pw.Deleted && !includeDeleted {
			continue
		}

		err = pw.decrypt(u)
		if err != nil {
			return nil, err
		}

		pw.UUID, err = uuid.Parse(uuidString)
		if err != nil {
			return nil, err
		}

		pw.UserID, err = uuid.Parse(useridString)
		if err != nil {
			return nil, err
		}

		pws = append(pws, pw)
	}

	pws.Sort()
	return pws, nil
}

func (m *PasswordModel) GetAllEncryptedForUser(u User) (PasswordList, error) {
	stmt := `SELECT id, uuid, userId, service, username, password, last_changed, deleted FROM passwords WHERE userId = ?`
	rows, err := m.DB.Query(stmt, u.ID.String())
	if err != nil {
		return nil, err
	}

	pws := PasswordList{}
	for rows.Next() {
		pw := Password{}
		var uuidString string
		var useridString string
		err := rows.Scan(&pw.ID, &uuidString, &useridString, &pw.EServiceName, &pw.EUsername, &pw.EPassword, &pw.LastChanged, &pw.Deleted)
		if err != nil {
			return nil, err
		}

		pw.UUID, err = uuid.Parse(uuidString)
		if err != nil {
			return nil, err
		}

		pw.UserID, err = uuid.Parse(useridString)
		if err != nil {
			return nil, err
		}

		pws = append(pws, pw)
	}

	return pws, nil
}

func (m *PasswordModel) DumbUpdate(p Password) error {
	stmt := `UPDATE passwords SET service = ?, username = ?, password = ?, last_changed = ?, deleted = ? WHERE uuid = ?`
	_, err := m.DB.Exec(stmt, p.EServiceName, p.EUsername, p.EPassword, p.LastChanged, p.Deleted, p.UUID.String())
	return err
}

func (m *PasswordModel) DumbInsert(p Password) error {
	stmt := `INSERT INTO passwords (uuid, userId, service, username, password, last_changed, deleted) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := m.DB.Exec(stmt, p.UUID.String(), p.UserID.String(), p.EServiceName, p.EUsername, p.EPassword, p.LastChanged, p.Deleted)
	return err
}

func (m *PasswordModel) Delete(UUID string) error {
	_, err := m.DB.Exec("DELETE FROM passwords WHERE uuid = ?", UUID)
	return err
}

func (m *PasswordModel) Exists(UUID string) (bool, error) {
	row := m.DB.QueryRow("SELECT EXISTS(SELECT uuid FROM passwords WHERE uuid = ?)", UUID)
	var result bool
	err := row.Scan(&result)
	return result, err
}

func (m *PasswordModel) ReplaceAllForUser(userID string, pwl PasswordList) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	for _, pw := range pwl {
		if !bytes.Equal(userUUID[:], pw.UserID[:]) {
			return errors.New("Passwords not for this user")
		}
	}

	_, err = m.DB.Exec(`DELETE FROM passwords WHERE userId = ?`, userID)
	if err != nil {
		return err
	}

	for _, pw := range pwl {
		err = m.DumbInsert(pw)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pl PasswordList) Search(searchTerm string) PasswordList {
	if strings.TrimSpace(searchTerm) == "" {
		return pl
	}

	searchTerm = strings.ToLower(searchTerm)

	res := PasswordList{}
	for _, p := range pl {
		if !p.isDecrypted() {
			continue
		}

		sn := strings.ToLower(p.ServiceName)
		if strings.Contains(sn, searchTerm) {
			res = append(res, p)
			continue
		}

		// un := strings.ToLower(p.Username)
		// if strings.Contains(un, searchTerm) {
		// 	res = append(res, p)
		// 	continue
		// }
	}
	return res
}

func (pl PasswordList) Sort() {
	sort.Slice(pl, func(i, j int) bool {
		return pl[i].ServiceName < pl[j].ServiceName
	})
}
