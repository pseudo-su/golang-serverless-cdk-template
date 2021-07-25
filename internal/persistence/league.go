package persistence

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	// "gorm.io/plugin/soft_delete"
)

// League schema
type League struct {
	// Gorm fields
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Fields
	Name        string `json:"name"`
	Description string `json:"description"`
}
