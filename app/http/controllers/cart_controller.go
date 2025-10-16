package controllers

import (
	"fmt"
	"strconv"
	"time"

	"goravel/app/models"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type CartController struct{}

func NewCartController() *CartController {
	return &CartController{}
}

func (c *CartController) Index(ctx http.Context) http.Response {
	type CartItemWithItem struct {
		ID          int     `json:"id"`
		ItemID      int     `json:"item_id"`
		Quantity    int     `json:"quantity"`
		HargaManual float64 `json:"harga_manual"`
		ItemName    string  `json:"item_name"`
		HargaBeli   float64 `json:"harga_beli"`
		Stok        int     `json:"stok"`
	}

	var results []CartItemWithItem

	query := `
	SELECT 
		c.id, c.item_id, c.quantity, c.harga_manual,
		COALESCE(i.nama_barang, '') AS item_name, 
		COALESCE(i.harga_beli, 0) AS harga_beli, 
		COALESCE(i.stok, 0) AS stok
	FROM cart_items c
	LEFT JOIN items i ON i.id = c.item_id
`

	err := facades.DB().Select(&results, query)
	if err != nil {
		return ctx.Response().Json(500, http.Json{"error": err.Error()})
	}

	return ctx.Response().Json(200, http.Json{"data": results})
}

// POST /cart/add/:id -> klik tombol + di FE
func (c *CartController) Add(ctx http.Context) http.Response {
    idParam := ctx.Request().Route("id")
    itemIDInt, err := strconv.Atoi(idParam)
    if err != nil {
        return ctx.Response().Json(400, http.Json{"error": "ID tidak valid"})
    }
    itemID := uint(itemIDInt)

    // quantity default = 1 (klik tombol + di FE)
    quantity := 1

    // cari item
    var item models.Item
    facades.Orm().Query().Where("id = ?", itemID).First(&item)

    // cek apakah item ditemukan
    if item.ID == 0 {
        return ctx.Response().Json(404, http.Json{"error": "Item tidak ditemukan"})
    }

    // cek stok
    if item.Stok < quantity {
        return ctx.Response().Json(400, http.Json{"error": "Stok tidak mencukupi"})
    }

    // update stok item
    item.Stok -= quantity
    if err := facades.Orm().Query().Save(&item); err != nil {
        return ctx.Response().Json(500, http.Json{"error": "Gagal update stok"})
    }

    // cek apakah item sudah ada di cart
    var cartItem models.CartItem
    facades.Orm().Query().Where("item_id = ?", itemID).First(&cartItem)

    if cartItem.ID != 0 {
        // sudah ada -> increment quantity
        cartItem.Quantity += quantity
        if err := facades.Orm().Query().Save(&cartItem); err != nil {
            return ctx.Response().Json(500, http.Json{"error": "Gagal update cart"})
        }
    } else {
        // belum ada -> insert baru
        newCart := models.CartItem{
            ItemID:      itemID,
            Quantity:    quantity,
            HargaManual: 0, // default dulu
        }
        if err := facades.Orm().Query().Create(&newCart); err != nil {
            return ctx.Response().Json(500, http.Json{"error": err.Error()})
        }
    }

    return ctx.Response().Json(200, http.Json{"message": "Item berhasil ditambahkan ke keranjang"})
}





// DELETE /cart/remove/:id
func (c *CartController) Remove(ctx http.Context) http.Response {
	idParam := ctx.Request().Route("id")
	cartID, _ := strconv.Atoi(idParam)

	var cartItem models.CartItem
	// Cari dulu apakah item ada
	if err := facades.Orm().Query().Find(&cartItem, cartID); err != nil || cartItem.ID == 0 {
		return ctx.Response().Json(404, http.Json{
			"error": "Cart item tidak ditemukan",
		})
	}

	// Hapus item (sama seperti Destroy di ItemController)
	result, err := facades.Orm().Query().Where("id = ?", cartID).Delete(&models.CartItem{})
	if err != nil {
		return ctx.Response().Json(500, http.Json{
			"error": err.Error(),
		})
	}

	if result.RowsAffected == 0 {
		return ctx.Response().Json(404, http.Json{
			"error": "Cart item tidak ditemukan",
		})
	}

	return ctx.Response().Json(200, http.Json{
		"message": "Barang di keranjang berhasil dihapus",
	})
}


// PATCH /cart/update-harga/:id
func (c *CartController) UpdateHarga(ctx http.Context) http.Response {
	idParam := ctx.Request().Route("id")
	cartID, err := strconv.Atoi(idParam)
	if err != nil {
		return ctx.Response().Json(400, http.Json{"error": "ID tidak valid"})
	}

	// Ambil input harga_manual sebagai float64
	hargaManual, err := strconv.ParseFloat(ctx.Request().Input("harga_manual"), 64)
	if err != nil {
		return ctx.Response().Json(400, http.Json{"error": "Harga harus berupa angka"})
	}

	if hargaManual < 0 {
		return ctx.Response().Json(400, http.Json{"error": "Harga tidak valid"})
	}

	var cartItem models.CartItem
	if err := facades.Orm().Query().With("Item").Find(&cartItem, cartID); err != nil || cartItem.ID == 0 {
    	return ctx.Response().Json(404, http.Json{"error": "Cart item tidak ditemukan"})
	}

	// Update harga_manual (langsung, tanpa pointer)
	if hargaManual < cartItem.Item.HargaBeli {
	return ctx.Response().Json(400, http.Json{
		"error": fmt.Sprintf("Harga manual tidak boleh lebih kecil dari harga beli (%.0f)", cartItem.Item.HargaBeli),
		})
	}

	cartItem.HargaManual = hargaManual

	if err := facades.Orm().Query().Save(&cartItem); err != nil {
		return ctx.Response().Json(500, http.Json{"error": "Gagal memperbarui harga"})
	}

	return ctx.Response().Json(200, http.Json{
		"message": "Harga berhasil diperbarui",
		"data":    cartItem,
	})
}


// POST /cart/checkout
func (c *CartController) Checkout(ctx http.Context) http.Response {
	type request struct {
		NamaPembeli string  `json:"nama_pembeli"`
		NoHp        string  `json:"no_hp"`
		Alamat      string  `json:"alamat"`
		UangTitipan float64 `json:"uang_titipan"`
	}
	var body request
	if err := ctx.Request().Bind(&body); err != nil {
		return ctx.Response().Json(400, http.Json{"error": err.Error()})
	}

	// ambil cart + relasi item
	var cartItems []models.CartItem
	if err := facades.Orm().Query().With("Item").Get(&cartItems); err != nil {
		return ctx.Response().Json(500, http.Json{"error": err.Error()})
	}
	if len(cartItems) == 0 {
		return ctx.Response().Json(400, http.Json{"error": "Keranjang kosong"})
	}

	// struct untuk balikin ke FE
	type NotaItem struct {
		NamaBarang  string  `json:"nama_barang"`
		Jumlah      int     `json:"jumlah"`
		HargaSatuan float64 `json:"harga_satuan"`
		TotalHarga  float64 `json:"total_harga"`
	}

	notaItems := []NotaItem{}
	var total float64

	for _, cart := range cartItems {
		if cart.HargaManual == 0 {
			return ctx.Response().Json(400, http.Json{"error": "Harga manual belum diisi"})
		}

		item := cart.Item
		if cart.HargaManual < item.HargaBeli {
			return ctx.Response().Json(400, http.Json{"error": "Harga jual tidak boleh lebih kecil dari harga kulak"})
		}

		trx := models.Transaction{
			ItemID:      cart.ItemID,
			NamaPembeli: body.NamaPembeli,
			NoHp:        body.NoHp,
			Alamat:      body.Alamat,
			Jumlah:      cart.Quantity,
			TotalHarga:  cart.HargaManual,
			HargaSatuan: cart.HargaManual,
			Tanggal:     time.Now(),
		}
		facades.Orm().Query().Create(&trx)

		// simpan versi "nota"
		notaItems = append(notaItems, NotaItem{
			NamaBarang:  item.NamaBarang,
			Jumlah:      cart.Quantity,
			HargaSatuan: cart.HargaManual,
			TotalHarga:  cart.HargaManual,
		})

		total += trx.TotalHarga
	}

	// kosongkan cart
	facades.Orm().Query().Exec("DELETE FROM cart_items")

	totalKurang := total - body.UangTitipan

	// balikin JSON siap pakai untuk FE
	return ctx.Response().Json(200, http.Json{
		"message": "Checkout berhasil",
		"nota": map[string]interface{}{
			"tanggal":      time.Now().Format("2006-01-02"),
			"nama_pembeli": body.NamaPembeli,
			"no_hp":        body.NoHp,
			"alamat":       body.Alamat,
			"uang_titipan": body.UangTitipan,
			"total":        total,
			"total_kurang": totalKurang,
			"items":        notaItems,
		},
	})
}

