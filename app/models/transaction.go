package models

import "time"

type Transaction struct {
	ID          uint      `gorm:"primaryKey"`
	ItemID      uint      // foreign key
	Jumlah      int
	NamaPembeli string    `gorm:"size:255"`
	NoHp        string    `gorm:"size:20"`
	Alamat      string    `gorm:"size:255"`
	TotalHarga  float64   `gorm:"type:decimal(12,2)"`
	HargaSatuan float64   `gorm:"type:decimal(12,2)"`
	Tanggal     time.Time

	// Relasi
	Item Item `gorm:"foreignKey:ItemID" json:"item"`
}
