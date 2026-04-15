package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel mirrors gorm.Model but with snake_case JSON tags
// so frontend can always rely on `id`, `created_at`, `updated_at`.
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
