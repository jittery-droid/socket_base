package models

import (
	"github.com/jinzhu/gorm"
)

type ServicesConfig func(*Services) error

type Services struct {
	User   UserService
	Friend FriendService
	// Socket SocketService
	db *gorm.DB
}

func WithGorm(dialect, dbInfo string) ServicesConfig {
	return func(s *Services) error {
		db, err := gorm.Open(dialect, dbInfo)
		if err != nil {
			return err
		}
		s.db = db
		return nil
	}
}

func WithLogMode(mode bool) ServicesConfig {
	return func(s *Services) error {
		s.db.LogMode(mode)
		return nil
	}
}

func WithUser(pepper, jwtSecret string) ServicesConfig {
	return func(s *Services) error {
		s.User = NewUserService(s.db, pepper, jwtSecret)
		return nil
	}
}

func WithFriend() ServicesConfig {
	return func(s *Services) error {
		s.Friend = NewFriendService(s.db)
		return nil
	}
}

func NewServices(cfgs ...ServicesConfig) (*Services, error) {
	var s Services
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

func (s *Services) Close() error {
	// also close socket service and user service?
	return s.db.Close()
}

// DestructiveReset drops all tables and rebuilds them
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}

func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Friend{}).Error
}
