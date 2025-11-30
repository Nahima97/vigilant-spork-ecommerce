package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Review struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Rating      int        `json:"rating"`
	ProductID   uuid.UUID  `json:"product_id"`
	UserID      uuid.UUID  `json:"user_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}
