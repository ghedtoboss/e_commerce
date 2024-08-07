package models

import "time"

type Order struct {
	ID          uint    `gorm:"primaryKey"`
	UserID      uint    `gorm:"not null"` // Siparişi veren kullanıcı
	User        User    `gorm:"foreignKey:UserID"`
	TotalAmount float64 `gorm:"not null"` // Toplam tutar
	Status      string  `gorm:"not null"` // Sipariş durumu: pending, confirmed, shipped, delivered, cancelled
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
