package models

import (
	"github.com/gofrs/uuid"
	"time"
)

type Cart struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID     uuid.UUID  `json:"user_id"`
	User       User       `gorm:"foreignKey:UserID" json:"user"`
	Items      []CartItem `json:"items"`
	TotalPrice float64    `json:"total_price"`
}

type CartItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CartID    uuid.UUID `json:"cart_id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	UnitPrice float64   `json:"unit_price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
