package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/go-jose/go-jose/v4"
	"sutext.github.io/entry/scope"
	"sutext.github.io/suid"
)

type UserToken struct {
	UserID       suid.SUID `json:"user_id"`
	ExpiryIn     time.Time `json:"expiry_in"`
	AccessToken  *string   `json:"access_token"`
	RefreshToken *string   `json:"refresh_token"`
}

// VerificationKey is a rotated signing key which can still be used to verify
// signatures.
type VerificationKey struct {
	PublicKey *jose.JSONWebKey `json:"publicKey"`
	Expiry    time.Time        `json:"expiry"`
}

// Keys hold encryption and signing keys.
type Keys struct {
	SigningKey       *jose.JSONWebKey  `json:"signingKey"`
	SigningKeyPub    *jose.JSONWebKey  `json:"signingKeyPub"`
	VerificationKeys []VerificationKey `json:"verificationKeys"`
	NextRotation     time.Time         `json:"nextRotation"`
}

func (k Keys) Value() (driver.Value, error) {
	return json.Marshal(k)
}
func (k *Keys) Scan(src any) error {
	data, ok := src.([]byte)
	if !ok {
		return nil
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, k)
}
func (k Keys) GormDataType() string {
	return "blob"
}

type KeyRecord struct {
	ID   string `gorm:"primary_key"`
	Keys Keys
}

// AuthRequest represents a OAuth2 client authorization request. It holds the state
// of a single auth flow up to the point that the user authorizes the client.
type AuthRequest struct {
	ID                  string
	ConnID              string
	ClientID            string
	Scopes              scope.Scopes
	Nonce               string
	State               string
	Expiry              time.Time
	Claims              Claims
	PKCE                PKCE
	LoggedIn            bool
	HMACKey             []byte
	RedirectURI         string
	ResponseTypes       Strings
	ForceApprovalPrompt bool
}

// RefreshToken is an OAuth2 refresh token which allows a client to request new
// tokens on the end user's behalf.
type RefreshToken struct {
	ID            string `gorm:"primary_key"`
	ClientID      string
	UserID        suid.SUID
	Token         string
	ObsoleteToken string
	CreatedAt     time.Time
	LastUsed      time.Time
	Claims        Claims
	Scopes        scope.Scopes
	Nonce         string
}

// AuthCode represents a code which can be exchanged for an OAuth2 token response.
//
// This value is created once an end user has authorized a client, the server has
// redirect the end user back to the client, but the client hasn't exchanged the
// code for an access_token and id_token.
type AuthCode struct {
	ID          string
	ClientID    string
	RedirectURI string
	Nonce       string
	Scopes      scope.Scopes
	Claims      Claims
	Expiry      time.Time
	PKCE        PKCE
}
