package repository

import (
	"vigilant-spork/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
		GetOrderHistory(userID int) ([]models.Order, error)
}

type OrderRepo struct {
	Db *gorm.DB
}

func (r *OrderRepo) GetOrderHistory(userID int) ([]models.Order, error) {
	var orders []models.Order

	tx := r.Db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return orders, nil
}