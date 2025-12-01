package repository

import (
	"errors"
	"gorm.io/gorm"
	"vigilant-spork/db"
	"vigilant-spork/models"
)

type UserRepository interface {
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
	AddTokenToBlacklist(token string) error
	IsTokenBlacklisted(token string) (bool, error)
}

type UserRepo struct {
	Db *gorm.DB
}

func (r *UserRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := db.Db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) CreateUser(user *models.User) error {
	return db.Db.Create(user).Error
}

func (r *UserRepo) AddTokenToBlacklist(token string) error {
	var entry models.BlacklistedToken
	entry.Token = token
	return db.Db.Create(&entry).Error
}

func (r *UserRepo) IsTokenBlacklisted(token string) (bool, error) {
	var entry models.BlacklistedToken
	err := db.Db.Where("token = ?", token).First(&entry).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
