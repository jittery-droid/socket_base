package models

import (
	"sockets/hash"

	"github.com/jinzhu/gorm"
)

const hmacSecretKey = "secret-hmac-key"

type User struct {
	gorm.Model
	Name string
	Email string `gorm:"not null; unique_index"`
	Password string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember string `gorm:"-"`
	RememberHash string `gorm:"not null; unique_index"`
}

type UserDB interface {
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByToken(token string) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(user *User) error
}

type userDbHandle struct {
	db *gorm.DB
}

type userValidator struct {
	UserDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
}

type userService struct {
	UserDB
}

func NewUserService(db *gorm.DB) UserService {
	ud := &userDbHandle{db}
	hmac := hash.NewHMAC(hmacSecretKey)
	uv := newUserValidationLayer(ud, hmac)
	return &userService{
		UserDB: uv,
	}
}

func newUserValidationLayer(udb UserDB, hmac hash.HMAC) *userValidator {
	return &userValidator{
		UserDB: udb,
		hmac: hmac,
		emailRegex: regexp.MustCompile(
			`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`)
		}
	}
}
