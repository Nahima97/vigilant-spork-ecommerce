package repository

import (
	"vigilant-spork/db"
	"vigilant-spork/models"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	AddProduct(product *models.Product) error
	GetProductByID(id uuid.UUID) (*models.Product, error)
	GetProductByName(name string) (*models.Product, error)
}

type ProductRepo struct {
	Db *gorm.DB
}

func (r *ProductRepo) AddProduct(product *models.Product) error {
	return db.Db.Create(product).Error
}

func (r *ProductRepo) GetProductByID(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := db.Db.First(&product, id).Error 
	if err != nil {
		return nil, err 
	}
	return &product, nil
}

func (r *ProductRepo) GetProductByName(name string) (*models.Product, error) {
	var product models.Product
	err := db.Db.Where("name = ?", name).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}