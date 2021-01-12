package models

import (
	"github.com/jinzhu/gorm"
)

type Friend struct {
	gorm.Model
	UserID   uint   `gorm:"not_null;index"`
	FriendID uint   `gorm:"not_null"`
	Status   string `gorm:"not_null"`
}

type FriendService interface {
	FriendDB
}

type FriendDB interface {
	ByID(id uint) (*Friend, error)
	ByUserID(userID uint) ([]Friend, error)
	Create(friend *Friend) error
	Update(friend *Friend) error
	Delete(id uint) error
}

type friendService struct {
	FriendDB
}

type friendValidator struct {
	FriendDB
}

type friendGorm struct {
	db *gorm.DB
}

type friendValFunc func(*Friend) error

func runFriendValFuncs(friend *Friend, fns ...friendValFunc) error {
	for _, fn := range fns {
		if err := fn(friend); err != nil {
			return err
		}
	}
	return nil
}

func NewFriendService(db *gorm.DB) FriendService {
	return &friendService{
		FriendDB: &friendValidator{&friendGorm{db}},
	}
}

func (fv *friendValidator) Create(friend *Friend) error {
	err := runFriendValFuncs(friend,
		fv.userIDRequired,
		fv.friendIDRequired)
	if err != nil {
		return err
	}
	return fv.FriendDB.Create(friend)
}

func (fv *friendValidator) Update(friend *Friend) error {
	err := runFriendValFuncs(friend,
		fv.userIDRequired,
		fv.friendIDRequired)
	if err != nil {
		return err
	}
	return fv.FriendDB.Update(friend)
}

func (fv *friendValidator) Delete(id uint) error {
	if id <= 0 {
		return ErrIDInvalid
	}
	return fv.FriendDB.Delete(id)
}

func (fv *friendValidator) userIDRequired(f *Friend) error {
	if f.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (fv *friendValidator) friendIDRequired(f *Friend) error {
	if f.FriendID <= 0 {
		return ErrFriendIDRequired
	}
	return nil
}

func (fg *friendGorm) ByID(id uint) (*Friend, error) {
	var friend Friend
	db := fg.db.Where("id = ?", id)
	err := first(db, &friend)
	return &friend, err
}

func (fg *friendGorm) ByUserID(userID uint) ([]Friend, error) {
	var friends []Friend
	err := fg.db.Where("user_id = ?", userID).Find(&friends).Error
	if err != nil {
		return nil, err
	}
	return friends, nil
}

func (fg *friendGorm) Create(friend *Friend) error {
	return fg.db.Create(friend).Error
}

func (fg *friendGorm) Update(friend *Friend) error {
	return fg.db.Save(friend).Error
}

func (fg *friendGorm) Delete(id uint) error {
	friend := Friend{Model: gorm.Model{ID: id}}
	return fg.db.Delete(&friend).Error
}
