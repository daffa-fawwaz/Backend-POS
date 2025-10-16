package routes

import (
	"goravel/app/http/controllers"

	"github.com/goravel/framework/facades"
)

func Api() {
	itemController := controllers.NewItemController()

	// CRUD Routes Items
	facades.Route().Get("/items", itemController.Index)           // List semua item
	facades.Route().Post("/items", itemController.Store)          // Tambah item baru
	facades.Route().Get("/items/{id}", itemController.Show)       // Detail item
	facades.Route().Put("/items/{id}", itemController.Update)     // Update item
	facades.Route().Delete("/items/{id}", itemController.Destroy) // Hapus item

	// Custom Routes
	facades.Route().Get("/items/search", itemController.Search) // Cari item

	// Controller transaksi
	transactionController := controllers.NewTransactionController()

	// Checkout transaksi (POST /transactions/checkout/:id)
	facades.Route().Post("/transactions/checkout/{id}", transactionController.ProcessCheckout)

	// Chart pendapatan bulanan (GET /transactions/chart)
	facades.Route().Get("/transactions/chart", transactionController.ChartPendapatanBulanan)

	// List semua transaksi (GET /transactions)
	facades.Route().Get("/transactions", transactionController.GetAll)


	// Detail transaksi (GET /transactions/:id)
	facades.Route().Get("/transactions/{id}", transactionController.GetByID)

	// Controller cart
	cartController := controllers.NewCartController()

	// Cart Routes
	facades.Route().Get("/cart", cartController.Index)                     // List semua item di cart
	facades.Route().Post("/cart/add/{id}", cartController.Add)             // Tambah item ke cart
	facades.Route().Delete("/cart/remove/{id}", cartController.Remove)     // Hapus item dari cart
	facades.Route().Patch("/cart/update-harga/{id}", cartController.UpdateHarga) // Update harga item
	facades.Route().Post("/cart/checkout", cartController.Checkout)        // Checkout cart
}

