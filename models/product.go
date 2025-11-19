package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Product struct {
    ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    Name          string    `json:"name"`
    Description   string    `gorm:"type:text" json:"description"`
    Category      string    `json:"category"`
    Price         float64   `json:"price"`
    StockQuantity int       `json:"stockQuantity"`
    Data          string    `json:"data"`
    CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	DeletedAt   *time.Time  `json:"deleted_at"`
}
