package models

import (
	"time"

	"gorm.io/gorm"
)

type Shop struct {
	ID        uint           `gorm:"primaryKey"`
	Name      string         `gorm:"not null"`
	OwnerID   uint           `gorm:"not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
