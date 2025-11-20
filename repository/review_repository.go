package repository

import"gorm.io/gorm"

type ReviewRepository interface {
}

type ReviewRepo struct {
	Db *gorm.DB
}