package models

import "time"

type Item struct {
	ID           uint           `gorm:"primaryKey"`
	NamaBarang   string         `gorm:"size:255"`
	TipeBarang   string         `gorm:"size:100"`
	HargaJual    float64        `gorm:"type:decimal(12,2)"`
	HargaBeli    float64        `gorm:"type:decimal(12,2)"`
	TanggalOrder time.Time
	Stok         int

	// Relasi
	Transactions []Transaction `gorm:"foreignKey:ItemID"`
	CartItems    []CartItem    `gorm:"foreignKey:ItemID"`
}
