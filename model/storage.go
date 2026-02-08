package model

import (
	"context"
	"crypto"
	"crypto/rand"
	"encoding/base32"
	"io"
	"strings"

	"gorm.io/gorm"
	"sutext.github.io/suid"
)

// Kubernetes only allows lower case letters for names.
//
// TODO(ericchiang): refactor ID creation onto the storage.
var encoding = base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567")

// Valid characters for user codes
// const validUserCharacters = "BCDFGHJKLMNPQRSTVWXZ"

// NewDeviceCode returns a 32 char alphanumeric cryptographically secure string
func NewDeviceCode() string {
	return newSecureID(32)
}

// NewID returns a random string which can be used as an ID for objects.
func NewID() string {
	return newSecureID(16)
}

func newSecureID(len int) string {
	buff := make([]byte, len) // random ID.
	if _, err := io.ReadFull(rand.Reader, buff); err != nil {
		panic(err)
	}
	// Avoid the identifier to begin with number and trim padding
	return string(buff[0]%26+'a') + strings.TrimRight(encoding.EncodeToString(buff[1:]), "=")
}

// NewHMACKey returns a random key which can be used in the computation of an HMAC
func NewHMACKey(h crypto.Hash) []byte {
	return []byte(newSecureID(h.Size()))
}

type Storage interface {
	GetKeys(ctx context.Context) (Keys, error)
	UpdateKeys(ctx context.Context, updater func(old Keys) (Keys, error)) error

	GetUser(ctx context.Context, id suid.SUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByPhone(ctx context.Context, phone string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error

	GetToken(ctx context.Context, id suid.SUID) (*UserToken, error)
	CreateToken(ctx context.Context, token *UserToken) error
	DeleteToken(ctx context.Context, token *UserToken) error

	GetClient(ctx context.Context, id string) (Client, error)
	CreateClient(ctx context.Context, client Client) error
	DeleteClient(ctx context.Context, id string) error
	UpdateClient(ctx context.Context, id string, updater func(c Client) (Client, error)) error
	ListClients(ctx context.Context) ([]Client, error)

	GetAuthRequest(ctx context.Context, id string) (AuthRequest, error)
	CreateAuthRequest(ctx context.Context, a AuthRequest) error
	DeleteAuthRequest(ctx context.Context, id string) error
	UpdateAuthRequest(ctx context.Context, id string, updater func(a AuthRequest) (AuthRequest, error)) error

	GetAuthCode(ctx context.Context, id string) (AuthCode, error)
	CreateAuthCode(ctx context.Context, c AuthCode) error
	DeleteAuthCode(ctx context.Context, id string) error
	UpdateAuthCode(ctx context.Context, id string, updater func(c AuthCode) (AuthCode, error)) error

	GetRefresh(ctx context.Context, id string) (RefreshToken, error)
	GetRefreshByUserAndClient(ctx context.Context, userID suid.SUID, clientID string) (RefreshToken, error)
	CreateRefresh(ctx context.Context, r RefreshToken) error
	DeleteRefresh(ctx context.Context, id string) error
	UpdateRefresh(ctx context.Context, id string, updater func(r RefreshToken) (RefreshToken, error)) error

	CreateTokenInfo(ctx context.Context) (*TokenInfo, error)
}
type Driver interface {
	Open() (db *gorm.DB, err error)
}

func Open(driver Driver) (Storage, error) {
	db, err := driver.Open()
	if err != nil {
		return nil, err
	}
	AutoMigrate(db)
	return &storage{
		db: db,
	}, nil
}

type storage struct {
	db *gorm.DB
}

func (s *storage) GetKeys(ctx context.Context) (Keys, error) {
	var record KeyRecord
	err := s.db.WithContext(ctx).First(&record).Error
	if err != nil {
		return Keys{}, err
	}
	return record.Keys, nil
}
func (s *storage) UpdateKeys(ctx context.Context, updater func(old Keys) (Keys, error)) error {
	var record KeyRecord
	err := s.db.WithContext(ctx).First(&record).Error
	if err != nil {
		return err
	}
	newKeys, err := updater(record.Keys)
	if err != nil {
		return err
	}
	record.Keys = newKeys
	return s.db.WithContext(ctx).Save(&record).Error
}

// Below is User implementations
func (s *storage) GetUser(ctx context.Context, id suid.SUID) (*User, error) {
	var user User
	err := s.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (s *storage) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := s.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (s *storage) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := s.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *storage) GetUserByPhone(ctx context.Context, phone string) (*User, error) {
	var user User
	err := s.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (s *storage) CreateUser(ctx context.Context, user *User) error {
	return s.db.WithContext(ctx).Create(user).Error
}

func (s *storage) UpdateUser(ctx context.Context, user *User) error {
	return s.db.WithContext(ctx).Save(user).Error
}

func (s *storage) GetToken(ctx context.Context, id suid.SUID) (*UserToken, error) {
	var token UserToken
	err := s.db.WithContext(ctx).First(&token, id).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (s *storage) CreateToken(ctx context.Context, token *UserToken) error {
	return s.db.WithContext(ctx).Create(token).Error
}

func (s *storage) DeleteToken(ctx context.Context, token *UserToken) error {
	return s.db.WithContext(ctx).Delete(token).Error
}

// Below is Client implementations
func (s *storage) GetClient(ctx context.Context, id string) (Client, error) {
	var client Client
	err := s.db.WithContext(ctx).First(&client, id).Error
	return client, err
}

func (s *storage) CreateClient(ctx context.Context, client Client) error {
	return s.db.WithContext(ctx).Create(client).Error
}

func (s *storage) UpdateClient(ctx context.Context, id string, updater func(c Client) (Client, error)) error {
	var client Client
	err := s.db.WithContext(ctx).First(&client, id).Error
	if err != nil {
		return err
	}
	newClient, err := updater(client)
	if err != nil {
		return err
	}
	return s.db.WithContext(ctx).Save(&newClient).Error
}
func (s *storage) DeleteClient(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Delete(Client{}, id).Error
}
func (s *storage) ListClients(ctx context.Context) ([]Client, error) {
	var clients []Client
	err := s.db.WithContext(ctx).Find(&clients).Error
	if err != nil {
		return nil, err
	}
	return clients, nil
}

// Below is AuthRequest implementations
func (s *storage) GetAuthRequest(ctx context.Context, id string) (AuthRequest, error) {
	var ar AuthRequest
	err := s.db.WithContext(ctx).First(&ar, id).Error
	if err != nil {
		return AuthRequest{}, err
	}
	return ar, nil
}

func (s *storage) CreateAuthRequest(ctx context.Context, a AuthRequest) error {
	return s.db.WithContext(ctx).Create(a).Error
}

func (s *storage) DeleteAuthRequest(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Delete(AuthRequest{}, id).Error
}

func (s *storage) UpdateAuthRequest(ctx context.Context, id string, updater func(a AuthRequest) (AuthRequest, error)) error {
	var a AuthRequest
	err := s.db.WithContext(ctx).First(&a, id).Error
	if err != nil {
		return err
	}
	newA, err := updater(a)
	if err != nil {
		return err
	}
	return s.db.WithContext(ctx).Save(&newA).Error
}

// Below is AuthCode implementations
func (s *storage) GetAuthCode(ctx context.Context, id string) (AuthCode, error) {
	var code AuthCode
	err := s.db.WithContext(ctx).First(&code, id).Error
	if err != nil {
		return AuthCode{}, err
	}
	return code, nil
}

func (s *storage) CreateAuthCode(ctx context.Context, c AuthCode) error {
	return s.db.WithContext(ctx).Create(c).Error
}

func (s *storage) DeleteAuthCode(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Delete(AuthCode{}, id).Error
}

func (s *storage) UpdateAuthCode(ctx context.Context, id string, updater func(c AuthCode) (AuthCode, error)) error {
	var code AuthCode
	err := s.db.WithContext(ctx).First(&code, id).Error
	if err != nil {
		return err
	}
	newCode, err := updater(code)
	if err != nil {
		return err
	}
	return s.db.WithContext(ctx).Save(&newCode).Error
}

// Below is RefreshToken implementations
func (s *storage) GetRefresh(ctx context.Context, id string) (RefreshToken, error) {
	var refresh RefreshToken
	err := s.db.WithContext(ctx).First(&refresh, "id = ?", id).Error
	if err != nil {
		return RefreshToken{}, err
	}
	return refresh, nil
}
func (s *storage) GetRefreshByUserAndClient(ctx context.Context, userID suid.SUID, clientID string) (RefreshToken, error) {
	var refresh RefreshToken
	err := s.db.WithContext(ctx).First(&refresh, "user_id = ? AND client_id = ?", userID, clientID).Error
	if err != nil {
		return RefreshToken{}, err
	}
	return refresh, nil
}
func (s *storage) CreateRefresh(ctx context.Context, r RefreshToken) error {
	return s.db.WithContext(ctx).Create(r).Error
}

func (s *storage) DeleteRefresh(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Delete(RefreshToken{}, "id = ?", id).Error
}

func (s *storage) UpdateRefresh(ctx context.Context, id string, updater func(r RefreshToken) (RefreshToken, error)) error {
	var refresh RefreshToken
	err := s.db.WithContext(ctx).First(&refresh, "id = ?", id).Error
	if err != nil {
		return err
	}
	newRefresh, err := updater(refresh)
	if err != nil {
		return err
	}
	return s.db.WithContext(ctx).Save(&newRefresh).Error
}

func (s *storage) CreateTokenInfo(ctx context.Context) (ti *TokenInfo, err error) {
	return ti, nil
}
