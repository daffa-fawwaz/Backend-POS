package models

type CartItem struct {
	ID       uint `gorm:"primaryKey"`
	ItemID   uint
	Quantity int

	// Relasi
	Item Item `gorm:"foreignKey:ItemID"`
}
