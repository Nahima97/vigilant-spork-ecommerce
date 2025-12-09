package repository

import (
	"gorm.io/gorm"
	"vigilant-spork/db"
	"vigilant-spork/models"
)

type UserRepository interface {
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
	AddTokenToBlacklist(token string) error
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
	err := db.Db.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) AddTokenToBlacklist(token string) error {
	var entry models.BlacklistedToken
	entry.Token = token
	err := db.Db.Create(&entry).Error
	if err != nil {
		return err
	}
	return nil
}
