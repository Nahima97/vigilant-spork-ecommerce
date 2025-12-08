package repository

import (
	"errors"
	"vigilant-spork/db"
	"vigilant-spork/models"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type CartRepository interface {
	GetOrCreateCart(userID uuid.UUID) (*models.Cart, error)
	AddItemToCart(productID, cartID uuid.UUID) error
	GetCartItems(cartID uuid.UUID) ([]models.CartItem, error)
	UpdateCartTotal(total int64, cartID uuid.UUID) error
	GetCartByUserID(userID uuid.UUID) (*models.Cart, error)
	GetCartItemsByCartID(cartID uuid.UUID) ([]models.CartItem, error)
}

type CartRepo struct {
	Db *gorm.DB
}

func (r *CartRepo) GetOrCreateCart(userID uuid.UUID) (*models.Cart, error) {
	var cart models.Cart
	err := db.Db.Preload("User").Preload("Items").Where("user_id = ?", userID).First(&cart).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		cart = models.Cart{
			UserID: userID,
			Total:  0,
		}
		err = db.Db.Create(&cart).Error
		if err != nil {
			return nil, err
		}

		err = db.Db.Preload("User").Preload("Items").First(&cart, cart.ID).Error
		if err != nil {
			return nil, err
		}

		return &cart, nil
	}
	if err != nil {
		return nil, err
	}

	return &cart, nil
}

func (r *CartRepo) AddItemToCart(productID, cartID uuid.UUID) error {
	var cart models.Cart
	err := db.Db.Where("id = ?", cartID).First(&cart).Error
	if err != nil {
		return err
	}

	var product models.Product
	err = db.Db.Where("id = ?", productID).First(&product).Error
	if err != nil {
		return err
	}

	var cartItem models.CartItem
	err = db.Db.Where("cart_id = ? AND product_id = ?", cart.ID, productID).First(&cartItem).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if product.StockQuantity < 1 {
			return ErrInsufficientStock
		}
		cartItem = models.CartItem{
			CartID:    cart.ID,
			ProductID: productID,
			Quantity:  1,
			UnitPrice: product.Price,
		}
		return db.Db.Create(&cartItem).Error
	}
	if err != nil {
		return err
	}

	if product.StockQuantity < cartItem.Quantity+1 {
		return ErrInsufficientStock
	}

	cartItem.Quantity++
	cartItem.UnitPrice = product.Price

	err = db.Db.Save(&cartItem).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *CartRepo) GetCartItems(cartID uuid.UUID) ([]models.CartItem, error) {
	var items []models.CartItem
	err := db.Db.Where("cart_id = ?", cartID).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *CartRepo) UpdateCartTotal(total int64, cartID uuid.UUID) error {
	err := db.Db.Model(&models.Cart{}).Where("id = ?", cartID).Update("total", total).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CartRepo) GetCartByUserID(userID uuid.UUID) (*models.Cart, error) {
	var cart models.Cart
	err := db.Db.Preload("Items.Product").Preload("User").Where("user_id = ?", userID).First(&cart).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *CartRepo) GetCartItemsByCartID(cartID uuid.UUID) ([]models.CartItem, error) {
	var items []models.CartItem
	err := db.Db.Preload("Product").Where("cart_id = ?", cartID).Find(&items).Error
	return items, err
}
