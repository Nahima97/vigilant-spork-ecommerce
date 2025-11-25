package repository

import (
	"vigilant-spork/db"
	"vigilant-spork/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
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

//SEVI

func (r *UserRepo) AddTokenToBlacklist(token string) error {
	entry := models.BlacklistedToken{
		Token: token,
	}
	return r.Db.Create(&entry).Error
}

func (r *UserRepo) IsTokenBlacklisted(token string) bool {
	var count int64
	r.Db.Model(&models.BlacklistedToken{}).
		Where("token = ?", token).
		Count(&count)
	return count > 0
}


