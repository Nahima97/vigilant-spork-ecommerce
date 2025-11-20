package repository

import"gorm.io/gorm"

type UserRepository interface {
}

type UserRepo struct {
	Db *gorm.DB
}