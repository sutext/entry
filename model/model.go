package model

import (
	"database/sql/driver"
	"encoding/json"
	"slices"

	"gorm.io/gorm"
	"sutext.github.io/suid"
)

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&KeyRecord{},
		&User{},
		&Client{},
		&UserToken{},
		&AuthRequest{},
		&RefreshToken{},
		&AuthCode{},
	)
}

type Strings []string

func (s Strings) Contains(v string) bool {
	return slices.Contains(s, v)
}

func (s *Strings) Scan(src any) error {
	data, ok := src.([]byte)
	if !ok {
		return nil
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, s)
}

func (s Strings) Value() (driver.Value, error) {
	return json.Marshal(s)
}
func (s *Strings) GormDataType() string {
	return "blob"
}

// PKCE is a container for the data needed to perform Proof Key for Code Exchange (RFC 7636) auth flow
type PKCE struct {
	CodeChallenge       string
	CodeChallengeMethod string
}

// Claims represents the ID Token claims supported by the server.
type Claims struct {
	UserID   suid.SUID
	Username string
	Email    string
	Groups   []string
}

func (p PKCE) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *PKCE) Scan(src any) error {
	data, ok := src.([]byte)
	if !ok {
		return nil
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, p)
}
func (p PKCE) GormDataType() string {
	return "blob"
}

func (c Claims) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *Claims) Scan(src any) error {
	data, ok := src.([]byte)
	if !ok {
		return nil
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, c)
}
func (c Claims) GormDataType() string {
	return "blob"
}
