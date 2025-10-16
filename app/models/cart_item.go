package models

import "time"

type CartItem struct {
	ID          uint      `gorm:"primaryKey;column:id" json:"id"`
	ItemID      uint      `gorm:"column:item_id" json:"item_id"`
	Quantity    int       `gorm:"column:quantity" json:"quantity"`
	HargaManual float64   `gorm:"column:harga_manual" json:"harga_manual"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`

	Item Item `gorm:"foreignKey:ItemID;references:ID" json:"item"`
}
