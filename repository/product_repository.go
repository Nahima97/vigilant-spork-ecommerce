package repository

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"vigilant-spork/db"
	"vigilant-spork/models"
)

type ProductRepository interface {
	AddProduct(product []models.Product) error
	GetProductByID(id uuid.UUID) (*models.Product, error)
	GetProductByName(name string) (*models.Product, error)
	GetProducts(limit int, offset int, minPrice int, maxPrice int, category string) ([]models.Product, error)
	GetProductsMetadata() (int64, error)
	UpdateProduct(product *models.Product) (*models.Product, error)
	DeleteProduct(id uuid.UUID) error
	UpdateAggregates(productID uuid.UUID, avgRating float64, reviewCount int64) error
}

type ProductRepo struct {
	Db *gorm.DB
}

func (r *ProductRepo) AddProduct(products []models.Product) error {
	for i := range products {
		err := db.Db.Create(&products[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ProductRepo) GetProductByID(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := db.Db.Preload("Reviews.User").First(&product, id).Error
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

func (r *ProductRepo) GetProducts(limit int, offset int, minPrice int, maxPrice int, category string) ([]models.Product, error) {
	var products []models.Product
	query := db.Db.Where("price BETWEEN ? AND ?", minPrice, maxPrice)

	if category != "" {
		query = query.Where("LOWER(category) = LOWER(?)", category)
	}
	err := query.Order("ID DESC").Limit(limit).Offset(offset).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepo) GetProductsMetadata() (int64, error) {
	var totalItems int64
	err := db.Db.Model(&models.Product{}).Count(&totalItems).Error
	if err != nil {
		return 0, err
	}
	return totalItems, nil
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
func (r *ProductRepo) UpdateAggregates(productID uuid.UUID, avgRating float64, reviewCount int64) error {
	return r.Db.Model(&models.Product{}).Where("id = ?", productID).
		Updates(map[string]interface{}{
			"rating":       avgRating,
			"review_count": reviewCount,
		}).Error
}
