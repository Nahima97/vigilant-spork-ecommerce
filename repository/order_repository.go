package repository

import (
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"vigilant-spork/models"
)

type OrderRepository interface {
	Transaction(ctx context.Context, fn func(repo OrderRepository) error) error
	GetCart(ctx context.Context, userID uuid.UUID) (*models.Cart, error)
	VerifyAndDeductStock(ctx context.Context, cartItem *models.CartItem) error
	CreateOrder(ctx context.Context, userID uuid.UUID) (*models.Order, error)
	MoveCartItemsToOrder(ctx context.Context, orderID uuid.UUID, cartID uuid.UUID) error
	ClearCart(ctx context.Context, cartID uuid.UUID) error
}

type OrderRepo struct {
	Db *gorm.DB
}

var ErrInsufficientStock = errors.New("insufficient stock")

// withTX creates a new repository instance with the given transaction
func (r *OrderRepo) withTX(tx *gorm.DB) *OrderRepo {
	return &OrderRepo{
		Db: tx,
	}
}

// Transaction manages the transaction lifecycle
func (r *OrderRepo) Transaction(ctx context.Context, fn func(repo OrderRepository) error) error {
	tx := r.Db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	repo := r.withTX(tx)
	err := fn(repo)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *OrderRepo) GetCart(ctx context.Context, userID uuid.UUID) (*models.Cart, error) {
	var cart models.Cart
	db := r.Db.WithContext(ctx)
	err := db.Where("user_id = ?", userID).First(&cart).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &cart, nil
}


func (r *OrderRepo) VerifyAndDeductStock(ctx context.Context, cartItem *models.CartItem) error {
	var product models.Product
	db := r.Db.WithContext(ctx)
	err := db.Model(&models.Product{}).Where("id = ?", cartItem.ProductID).Clauses(clause.Locking{Strength: "UPDATE"}).First(&product).Error
	if err != nil {
		return err
	}
	if cartItem.Quantity > product.StockQuantity {
		return ErrInsufficientStock
	}

	product.StockQuantity = product.StockQuantity - cartItem.Quantity

	err = db.Save(&product).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepo) CreateOrder(ctx context.Context, userID uuid.UUID) (*models.Order, error) {
	var order models.Order
	order.UserID = userID
	db := r.Db.WithContext(ctx)
	err := db.Create(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepo) MoveCartItemsToOrder(ctx context.Context, orderID uuid.UUID, cartID uuid.UUID) error {
	var cartItems []models.CartItem
	db := r.Db.WithContext(ctx)
	err := db.Where("cart_id = ?", cartID).Find(&cartItems).Error
	if err != nil {
		return err
	}

	if len(cartItems) == 0 {
		return nil
	}

	var orderItems []models.OrderItem

	for _, cartItem := range cartItems {
		orderItems = append(orderItems, models.OrderItem{
			OrderID:   orderID,
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
			UnitPrice: cartItem.UnitPrice,
		})
	}

	err = db.Create(&orderItems).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepo) ClearCart(ctx context.Context, cartID uuid.UUID) error {
	db := r.Db.WithContext(ctx)
	err := db.Where("cart_id = ?", cartID).Delete(&models.CartItem{}).Error
	if err != nil {
		return err
	}
	return nil
}
