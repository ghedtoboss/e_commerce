package models

import "time"

type OrderItem struct {
	ID        uint    `gorm:"primaryKey"`
	OrderID   uint    `gorm:"not null"` // Bağlı olduğu sipariş
	ProductID uint    `gorm:"not null"` // Ürün kimliği
	Quantity  int     `gorm:"not null"` // Miktar
	Price     float64 `gorm:"not null"` // Birim fiyat
	Total     float64 `gorm:"not null"` // Toplam fiyat
	CreatedAt time.Time
	UpdatedAt time.Time
}
