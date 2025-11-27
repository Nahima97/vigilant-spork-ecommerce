package services

import (
	"errors"
	"fmt"
	"vigilant-spork/models"
	"vigilant-spork/repository"

	"gorm.io/gorm"
)

type ProductService struct {
    ProductRepo repository.ProductRepository
}

func (s *ProductService) AddProduct (product *models.Product) error {
    if product.Name == "" {
        return errors.New("product name is required")
    }

    if product.Description == "" {
        return errors.New("product description is required")
    }

    if product.Category == "" {
        return errors.New("product category is required")
    }

    if product.Price == 0 {
        return errors.New("product price is required and cannot be 0")
    }

    if product.StockQuantity == 0 {
        return errors.New("stock quantity is required and cannot be 0")
    }

    existing, err := s.ProductRepo.GetProductByName(product.Name)
    if err == nil && existing != nil {
        return errors.New("product already exists")
    }

    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
    return fmt.Errorf("failed to check product existence: %w", err)
    }

	err = s.ProductRepo.AddProduct(product) 
    if err != nil {
        return err
    }
    return nil
}
