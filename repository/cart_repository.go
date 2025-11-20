package repository

import"gorm.io/gorm"

type CartRepository interface {
}

type CartRepo struct {
	Db *gorm.DB
}