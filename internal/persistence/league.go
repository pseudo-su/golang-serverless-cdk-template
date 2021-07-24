package persistence

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	// "gorm.io/plugin/soft_delete"
)

// League schema
type League struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	// IsDeleted   soft_delete.DeletedAt `gorm:"softDelete:flag"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
}
