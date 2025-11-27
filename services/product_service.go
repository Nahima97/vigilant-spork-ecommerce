package services

import (
	"errors"
	"fmt"
	"vigilant-spork/models"
	"vigilant-spork/repository"
    "github.com/google/uuid"
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

    existingProduct, err := s.ProductRepo.GetProductByName(product.Name)
    if err == nil && existingProduct != nil {
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

func (s *ProductService) GetProductByID (productID uuid.UUID) (*models.Product, error) {

    product, err := s.ProductRepo.GetProductByID(productID)
    if err != nil {
        return nil, err 
    }
    return product, nil
}

func (s *ProductService) GetProducts (page int, limit int) ([]models.Product, error) {
    offset := (page - 1) * limit 

    products, err := s.ProductRepo.GetProducts(limit, offset)
    if err != nil {
        return nil, err 
    }
    return products, nil
}

func (s *ProductService) UpdateProduct (productID uuid.UUID, req *models.Product) (*models.Product, error) {

    product, err := s.ProductRepo.GetProductByID(productID)
    if err != nil {
        return nil, err 
    }
    
    if req.Name != "" {
        product.Name = req.Name
    }
    if req.Description != "" {
        product.Description = req.Description
    }
    if req.Category != "" {
        product.Category = req.Category
    }
    if req.Price != 0.0 {
        product.Price = req.Price
    }
    if req.StockQuantity != 0 {
        product.StockQuantity = req.StockQuantity
    }

    updatedProduct, err := s.ProductRepo.UpdateProduct(product)
    if err != nil {
        return nil, err
    }
    return updatedProduct, nil
}

func (s *ProductService) DeleteProduct(productID uuid.UUID) error {

    err := s.ProductRepo.DeleteProduct(productID)
    if err != nil {
        return err 
    }
    return nil
}