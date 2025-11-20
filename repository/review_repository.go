package repository

import"gorm.io/gorm"

type ReviewRepository interface {
}

type ReviewRepo struct {
	db *gorm.DB
}