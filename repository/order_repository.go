package repository

import"gorm.io/gorm"

type OrderRepository interface {
}

type OrderRepo struct {
	db *gorm.DB
}