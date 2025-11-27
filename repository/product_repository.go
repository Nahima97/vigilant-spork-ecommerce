package repository

import (
	"vigilant-spork/db"
	"vigilant-spork/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	AddProduct(product *models.Product) error
	GetProductByID(id uuid.UUID) (*models.Product, error)
	GetProductByName(name string) (*models.Product, error)
	GetProducts(limit int, offset int) ([]models.Product, error)
	UpdateProduct(product *models.Product) (*models.Product, error)
	DeleteProduct(id uuid.UUID) error
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

// func (r *ProductRepo)

func (r *ProductRepo) GetProductByName(name string) (*models.Product, error) {
	var product models.Product
	err := db.Db.Where("name = ?", name).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepo) GetProducts(limit int, offset int) ([]models.Product, error) {
	var products []models.Product
	err := db.Db.Order("ID DESC").Limit(limit).Offset(offset).Find(&products).Error
	if err != nil {
		return nil, err 
	}
	return products, nil 
}

func (r *ProductRepo) UpdateProduct(product *models.Product) (*models.Product, error) {

	err := db.Db.Model(&models.Product{}).Where("id = ?", product.ID).Updates(product).Error
	if err != nil {
		return nil, err
	}

	var updatedProduct models.Product
	err = db.Db.First(&updatedProduct, "id = ?", product.ID).Error
	if err != nil {
		return nil, err
	}

	return &updatedProduct, nil
}

func (r *ProductRepo) DeleteProduct(id uuid.UUID) error {
	var product models.Product

	err := db.Db.Where("id = ?", id).Delete(&product).Error
	if err != nil {
		return err 
	}
	return nil
}