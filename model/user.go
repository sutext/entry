package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"sutext.github.io/suid"
)

type User struct {
	ID        suid.SUID  `json:"id" gorm:"primary_key" `
	Hash      string     `json:"hash"`
	Email     *string    `json:"email" gorm:"unique_index"`
	Phone     *string    `json:"phone" gorm:"unique_index"`
	Weight    uint       `json:"weight"` //g
	Height    uint       `json:"height"` //cm
	Avatar    *string    `json:"avatar"`
	Gender    Gender     `json:"gender" gorm:"type:tinyint(1)"`
	Groups    Strings    `json:"groups"`
	Username  *string    `json:"username" gorm:"unique_index"`
	Nickname  *string    `json:"nickname"`
	Birthday  *time.Time `json:"birthday"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
type UserView struct {
	ID       suid.SUID `json:"id"`
	Age      *int8     `json:"age,omitempty"`
	Email    *string   `json:"email,omitempty"`
	Phone    *string   `json:"phone,omitempty"`
	Weight   *uint     `json:"weight,omitempty"`
	Height   *uint     `json:"height,omitempty"`
	Avatar   *string   `json:"avatar,omitempty"`
	Gender   *Gender   `json:"gender,omitempty"`
	Username *string   `json:"username,omitempty"`
	Nickname *string   `json:"nickname,omitempty"`
	Birthday *int64    `json:"birthday,omitempty"`
}

func (u *User) VerifyPassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Hash), []byte(password)); err != nil {
		return false
	}
	return true
}
func (u *User) ToView() *UserView {
	var age *int8 = nil
	if u.Birthday != nil {
		v := int8(time.Now().Year() - u.Birthday.Year())
		age = &v
	}
	var birthday *int64 = nil
	if u.Birthday != nil {
		v := u.Birthday.Unix()
		birthday = &v
	}
	return &UserView{
		ID:       u.ID,
		Age:      age,
		Email:    u.Email,
		Phone:    u.Phone,
		Weight:   &u.Weight,
		Height:   &u.Height,
		Avatar:   u.Avatar,
		Gender:   &u.Gender,
		Username: u.Username,
		Nickname: u.Nickname,
		Birthday: birthday,
	}
}
func NewUser() *User {
	return &User{
		ID:        suid.New(),
		Gender:    GenderUnknown,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
func (u *User) Update(user *UserView) {
	if user.Email != nil {
		u.Email = user.Email
	}
	if user.Phone != nil {
		u.Phone = user.Phone
	}
	if user.Weight != nil {
		u.Weight = *user.Weight
	}
	if user.Height != nil {
		u.Height = *user.Height
	}
	if user.Avatar != nil {
		u.Avatar = user.Avatar
	}
	if user.Gender != nil {
		u.Gender = *user.Gender
	}
	if user.Username != nil {
		u.Username = user.Username
	}
	if user.Nickname != nil {
		u.Nickname = user.Nickname
	}
	if user.Birthday != nil {
		v := time.Unix(*user.Birthday, 0)
		u.Birthday = &v
	}
	u.UpdatedAt = time.Now()
}

type Gender int8

const (
	GenderUnknown Gender = iota
	GenderMale
	GenderFemale
)

func (g Gender) String() string {
	switch g {
	case GenderMale:
		return "male"
	case GenderFemale:
		return "female"
	default:
		return "unknown"
	}
}
