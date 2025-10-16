package controllers

import (
	"strconv"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"goravel/app/models"
)

type TransactionController struct{}

func NewTransactionController() *TransactionController {
	return &TransactionController{}
}

// POST /transactions/checkout/:id
func (c *TransactionController) ProcessCheckout(ctx http.Context) http.Response {
	idParam := ctx.Request().Route("id")
	itemID, _ := strconv.Atoi(idParam)

	type request struct {
		Jumlah      int     `json:"jumlah_beli"`
		TotalHarga  float64 `json:"total_harga"`
		NamaPembeli string  `json:"nama_pembeli"`
		NoHp        string  `json:"no_hp"`
		Alamat      string  `json:"alamat"`
	}
	var body request
	if err := ctx.Request().Bind(&body); err != nil {
		return ctx.Response().Json(400, http.Json{"error": err.Error()})
	}

	// Cari item
	var item models.Item
	if err := facades.Orm().Query().Where("id = ?", itemID).First(&item); err != nil {
		return ctx.Response().Json(404, http.Json{"error": "Item not found"})
	}

	// Validasi stok
	if body.Jumlah > item.Stok {
		return ctx.Response().Json(400, http.Json{"error": "Jumlah beli melebihi stok"})
	}

	// Hitung harga satuan
	hargaSatuan := body.TotalHarga / float64(body.Jumlah)

	// Tanggal otomatis hari ini
	tanggal := time.Now()

	// Simpan transaksi
	trx := models.Transaction{
		ItemID:      uint(itemID),
		Jumlah:      body.Jumlah,
		TotalHarga:  body.TotalHarga,
		HargaSatuan: hargaSatuan,
		Tanggal:     tanggal,
		NamaPembeli: body.NamaPembeli,
		NoHp:        body.NoHp,
		Alamat:      body.Alamat,
	}

	if err := facades.Orm().Query().Create(&trx); err != nil {
		return ctx.Response().Json(500, http.Json{"error": err.Error()})
	}

	// Kurangi stok
	item.Stok -= body.Jumlah
	if err := facades.Orm().Query().Save(&item); err != nil {
		return ctx.Response().Json(500, http.Json{"error": err.Error()})
	}

	// Load item ke relasi transaction
	var freshItem models.Item
	if err := facades.Orm().Query().Where("id = ?", trx.ItemID).First(&freshItem); err == nil {
		trx.Item = freshItem
	}

	return ctx.Response().Json(201, http.Json{
		"message": "Transaksi berhasil",
		"data":    trx,
	})
}


// GET /transactions/chart
func (c *TransactionController) ChartPendapatanBulanan(ctx http.Context) http.Response {
	now := time.Now()
	month := int(now.Month())
	year := now.Year()

	var transactions []models.Transaction
	if err := facades.Orm().Query().
		Where("EXTRACT(MONTH FROM tanggal) = ? AND EXTRACT(YEAR FROM tanggal) = ?", month, year).
		Get(&transactions); err != nil {
		return ctx.Response().Json(500, http.Json{"error": err.Error()})
	}

	// Load Item manual biar bisa hitung harga beli
	for i := range transactions {
		var item models.Item
		facades.Orm().Query().Where("id = ?", transactions[i].ItemID).First(&item)
		transactions[i].Item = item
	}

	chart := map[string]float64{
		"Minggu 1": 0,
		"Minggu 2": 0,
		"Minggu 3": 0,
		"Minggu 4": 0,
		"Minggu 5": 0,
	}

	for _, trx := range transactions {
		week := (trx.Tanggal.Day()-1)/7 + 1
		label := "Minggu " + strconv.Itoa(week)

		// Hitung keuntungan per transaksi
		keuntungan := (trx.HargaSatuan - trx.Item.HargaBeli) * float64(trx.Jumlah)
		chart[label] += keuntungan
	}

	return ctx.Response().Json(200, http.Json{"chart": chart})
}

// GET /transactions
func (c *TransactionController) GetAll(ctx http.Context) http.Response {
	type TransactionResponse struct {
		ID          uint      `json:"id"`
		NamaPembeli string    `json:"nama_pembeli"`
		Alamat      string    `json:"alamat"`
		Tanggal     time.Time `json:"tanggal"`
		NamaBarang  string    `json:"nama_barang"`
		Jumlah      int       `json:"jumlah"`
		TotalHarga  float64   `json:"total_harga"`
	}

	var results []TransactionResponse

	query := `
		SELECT 
			t.id,
			t.nama_pembeli,
			t.alamat,
			t.tanggal,
			i.nama_barang,
			t.jumlah,
			t.total_harga
		FROM transactions t
		LEFT JOIN items i ON i.id = t.item_id
		ORDER BY t.tanggal DESC
	`

	err := facades.DB().Select(&results, query)
	if err != nil {
		return ctx.Response().Json(500, http.Json{"error": err.Error()})
	}

	return ctx.Response().Json(200, http.Json{
		"message": "Berhasil mengambil semua transaksi",
		"data":    results,
	})
}


// GET /transactions/:id
func (c *TransactionController) GetByID(ctx http.Context) http.Response {
	idParam := ctx.Request().Route("id")
	trxID, err := strconv.Atoi(idParam)
	if err != nil {
		return ctx.Response().Json(400, http.Json{"error": "ID tidak valid"})
	}

	type TransactionDetail struct {
		ID          uint      `json:"id"`
		NamaPembeli string    `json:"nama_pembeli"`
		NoHp        string    `json:"no_hp"`
		Alamat      string    `json:"alamat"`
		Tanggal     time.Time `json:"tanggal"`
		NamaBarang  string    `json:"nama_barang"`
		Jumlah      int       `json:"jumlah"`
		HargaSatuan float64   `json:"harga_satuan"`
		TotalHarga  float64   `json:"total_harga"`
		HargaBeli   float64   `json:"harga_beli"`
		HargaJual   float64   `json:"harga_jual"`
	}

	var result TransactionDetail

	query := `
		SELECT 
			t.id,
			t.nama_pembeli,
			t.no_hp,
			t.alamat,
			t.tanggal,
			i.nama_barang,
			t.jumlah,
			t.harga_satuan,
			t.total_harga,
			i.harga_beli,
			i.harga_jual
		FROM transactions t
		LEFT JOIN items i ON i.id = t.item_id
		WHERE t.id = ?
		LIMIT 1
	`

	// ✅ gunakan Raw().Scan()
	if err := facades.Orm().Query().Raw(query, trxID).Scan(&result); err != nil {
		return ctx.Response().Json(500, http.Json{"error": err.Error()})
	}

	// Cek kalau datanya kosong (id gak ditemukan)
	if result.ID == 0 {
		return ctx.Response().Json(404, http.Json{"error": "Transaksi tidak ditemukan"})
	}

	return ctx.Response().Json(200, http.Json{
		"message": "Berhasil mengambil detail transaksi",
		"data":    result,
	})
}




