package repository

import (
	"vigilant-spork/models"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	UpdateAggregates(productID uuid.UUID, avgRating float64, reviewCount int64) error
}

type ProductRepo struct {
	Db *gorm.DB
}

func (r *ProductRepo) UpdateAggregates(productID uuid.UUID, avgRating float64, reviewCount int64) error {
	return r.Db.Model(&models.Product{}).Where("id = ?", productID).
		Updates(map[string]interface{}{
			"rating":       avgRating,
			"review_count": reviewCount,
		}).Error
}
