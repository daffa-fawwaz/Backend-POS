package routes

import (
	"github.com/goravel/framework/facades"
	"pos-api/app/http/controllers"
)

func Api() {
	// Transaction routes
	facades.Route().Get("/", controllers.TransactionController{}.ChartPendapatanBulanan)

	// Item routes
	facades.Route().Get("/items/search", controllers.ItemController{}.Search)
	facades.Route().Resource("items", controllers.ItemController{})
	facades.Route().Patch("/items/:id/update-stok", controllers.ItemController{}.UpdateStok)

	// Checkout routes
	facades.Route().Get("/items/:id/checkout", controllers.TransactionController{}.CheckoutForm)
	facades.Route().Post("/items/:id/checkout", controllers.TransactionController{}.ProcessCheckout)

	// History
	facades.Route().Get("/history", controllers.HistoryController{}.Index)

	// Nota
	facades.Route().Get("/nota/:id/cetak", controllers.TransactionController{}.CetakNota)

	// Cart routes
	facades.Route().Get("/cart", controllers.CartController{}.Index)
	facades.Route().Post("/cart/add/:id", controllers.CartController{}.Add)
	facades.Route().Post("/cart/remove/:id", controllers.CartController{}.Remove)
	facades.Route().Post("/cart/checkout", controllers.CartController{}.Checkout)
	facades.Route().Get("/cart/nota/:id", controllers.CartController{}.PrintNota)
	facades.Route().Put("/cart/update-harga/:id", controllers.CartController{}.UpdateHarga)
}
