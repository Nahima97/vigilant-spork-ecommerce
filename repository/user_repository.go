package repository

import (
	"vigilant-spork/db"
	"vigilant-spork/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByUsername(username string) (*models.User, error)
}

type UserRepo struct {
	Db *gorm.DB
}

func (r *UserRepo) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := db.Db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return &models.User{}, err
	}
	return &user, nil
}