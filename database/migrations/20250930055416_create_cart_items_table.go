package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250930092000CreateCartItemsTable struct{}

func (m *M20250930092000CreateCartItemsTable) Signature() string {
	return "2025_09_30_092000_create_cart_items_table"
}

func (m *M20250930092000CreateCartItemsTable) Up() error {
	return facades.Schema().Create("cart_items", func(table schema.Blueprint) {
		table.ID()
		table.UnsignedBigInteger("item_id")
		table.Integer("harga_manual").Nullable()
		table.Integer("quantity").Default(1)
		table.Timestamps()
	})
}

func (m *M20250930092000CreateCartItemsTable) Down() error {
	return facades.Schema().DropIfExists("cart_items")
}
