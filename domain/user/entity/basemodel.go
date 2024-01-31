package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel contains common fields for all models.
type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreatedAt time.Time
	CreatedBy uuid.UUID `gorm:"type:uuid;not null"`
	UpdatedAt time.Time
	UpdatedBy uuid.UUID
	DeletedAt gorm.DeletedAt `gorm:"index"`
	DeletedBy uuid.UUID
}
