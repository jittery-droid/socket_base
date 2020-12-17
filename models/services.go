package models

import (
	"github.com/jinzhu/gorm"
)

func InitServices(dbInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &Services{
		User:   NewUserService(db),
		Socket: NewSocketService(),
		db:     db,
	}, nil
}

type Services struct {
	User   UserService
	Socket SocketService
	db     *gorm.DB
}

func (s *Services) Close() error {
	// also close socket service and user service?
	return s.db.Close()
}

// DestructiveReset drops all tables and rebuilds them
func (s *Services) DestructiveReset() error {
	// drop other tables
	err := s.db.DropTableIfExists(&User{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}
