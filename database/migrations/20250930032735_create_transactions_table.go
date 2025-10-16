package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250930091000CreateTransactionsTable struct{}

func (m *M20250930091000CreateTransactionsTable) Signature() string {
	return "2025_09_30_091000_create_transactions_table"
}

func (m *M20250930091000CreateTransactionsTable) Up() error {
	return facades.Schema().Create("transactions", func(table schema.Blueprint) {
		table.ID()
		table.UnsignedBigInteger("item_id").Nullable()
		table.String("nama_pembeli", 255).Nullable()
		table.String("no_hp", 50).Nullable()
		table.String("alamat", 255).Nullable()
		table.Integer("jumlah")
		table.Double("total_harga")
		table.Double("harga_satuan").Nullable()
		table.Date("tanggal").Nullable()
		table.Timestamps()
	})
}

func (m *M20250930091000CreateTransactionsTable) Down() error {
	return facades.Schema().DropIfExists("transactions")
}
