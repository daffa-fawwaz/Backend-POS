package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250930023216CreateItemsTable struct{}

// Signature migration
func (m *M20250930023216CreateItemsTable) Signature() string {
	return "2025_09_30_023216_create_items_table"
}

// Up migration
func (m *M20250930023216CreateItemsTable) Up() error {
	return facades.Schema().Create("items", func(table schema.Blueprint) {
		table.ID()
		table.String("nama_barang", 255)
		table.String("tipe_barang", 100)
		table.Integer("stok")
		table.Decimal("harga_beli")
		table.Decimal("harga_jual").Nullable()
		table.Date("tanggal_order")
		table.Timestamps()
	})
}

// Down migration
func (m *M20250930023216CreateItemsTable) Down() error {
	return facades.Schema().DropIfExists("items")
}
