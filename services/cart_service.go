package services

import (
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"math"
	"vigilant-spork/models"
	"vigilant-spork/repository"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type CartService struct {
	CartRepo    repository.CartRepository
	ProductRepo repository.ProductRepository
}

func (s *CartService) AddToCart(userID, productID uuid.UUID) error {
	cart, err := s.CartRepo.GetOrCreateCart(userID)
	if err != nil {
		return err
	}

	if cart == nil {
		return fmt.Errorf("cart not found or could not be created for user")
	}

	product, err := s.ProductRepo.GetProductByID(productID)
	if err != nil {
		return err
	}

	if product == nil {
		return fmt.Errorf("product not found")
	}

	product.Price = int64(math.Round(float64(product.Price) * 100))

	err = s.CartRepo.AddItemToCart(productID, cart.ID)
	if err != nil {
		return err
	}

	items, err := s.CartRepo.GetCartItems(cart.ID)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return fmt.Errorf("items not found")
	}

	total := int64(0)
	for _, item := range items {
		total += int64(item.Quantity) * item.UnitPrice
	}

	err = s.CartRepo.UpdateCartTotal(total, cart.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *CartService) ViewCart(userID uuid.UUID) (*models.Cart, error) {
	cart, err := s.CartRepo.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}

	var total int64
	for i := range cart.Items {
		item := &cart.Items[i]
		item.UnitPrice = item.Product.Price
		total += int64(item.Quantity) * item.UnitPrice
	}
	cart.Total = total

	return cart, nil
}

func (s *CartService) UpdateItemQuantity(userID, productID uuid.UUID, quantity int) (*models.CartItem, error) {
	cartItem, err := s.CartRepo.UpdateItemQuantity(userID, productID, quantity)
	if err != nil {
		return nil, err 
	}
	return cartItem, nil 
func (s *CartService) RemoveItem(userID, productID uuid.UUID) error {
	cart, err := s.CartRepo.GetOrCreateCart(userID)
	if err != nil {
		return err
	}

	err = s.CartRepo.RemoveItemFromCart(cart.ID, productID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return gorm.ErrRecordNotFound
	}
	if err != nil {
		return err
	}

	items, err := s.CartRepo.GetCartItems(cart.ID)
	if err != nil {
		return err
	}

	total := int64(0)
	for _, item := range items {
		total += int64(item.Quantity) * item.UnitPrice
	}

	return s.CartRepo.UpdateCartTotal(total, cart.ID)
}
