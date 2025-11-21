package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string     `json:"name"`
	Email     string     `gorm:"unique"`
	Password  string     `json:"password"`
	Role      string     `json:"role"`
	CartID    uuid.UUID  `json:"cart_id"`
	OrderID   uuid.UUID  `json:"order_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}