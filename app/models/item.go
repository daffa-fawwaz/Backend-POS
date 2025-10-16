package models

import "time"

type Item struct {
	ID           uint      `gorm:"primaryKey;column:id" json:"id"`
	NamaBarang   string    `gorm:"size:255;column:nama_barang" json:"nama_barang"`
	TipeBarang   string    `gorm:"size:100;column:tipe_barang" json:"tipe_barang"`
	HargaJual    float64   `gorm:"type:decimal(18,2);column:harga_jual" json:"harga_jual"`
	HargaBeli    float64   `gorm:"type:decimal(18,2);column:harga_beli" json:"harga_beli"`
	TanggalOrder time.Time `gorm:"column:tanggal_order" json:"tanggal_order"`
	Stok         int       `gorm:"column:stok" json:"stok"`

	// Relasi
	Transactions []Transaction `gorm:"foreignKey:ItemID" json:"transactions,omitempty"`
	CartItems    []CartItem    `gorm:"foreignKey:ItemID" json:"cart_items,omitempty"`
}

