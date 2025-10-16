package controllers

import (
	"goravel/app/models"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type ItemController struct{}

func NewItemController() *ItemController {
	return &ItemController{}
}

// GET /items
func (c *ItemController) Index(ctx http.Context) http.Response {
	var items []models.Item
	if err := facades.Orm().Query().Get(&items); err != nil {
		return ctx.Response().Json(500, http.Json{"error": err.Error()})
	}
	return ctx.Response().Json(200, items)
}

// GET /items/{id}
func (c *ItemController) Show(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var item models.Item
	if err := facades.Orm().Query().Where("id = ?", id).First(&item); err != nil {
		return ctx.Response().Json(404, http.Json{"error": "Item not found"})
	}
	return ctx.Response().Json(200, item)
}

// POST /items
func (c *ItemController) Store(ctx http.Context) http.Response {
	var input models.Item
	if err := ctx.Request().Bind(&input); err != nil {
		return ctx.Response().Json(400, http.Json{"error": "Invalid input"})
	}

	if err := facades.Orm().Query().Create(&input); err != nil {
		return ctx.Response().Json(500, http.Json{"error": err.Error()})
	}
	return ctx.Response().Json(201, input)
}

// PUT /items/{id}
func (c *ItemController) Update(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	var item models.Item
	if err := facades.Orm().Query().Where("id = ?", id).First(&item); err != nil {
		return ctx.Response().Json(404, http.Json{"error": "Item not found"})
	}

	if err := ctx.Request().Bind(&item); err != nil {
		return ctx.Response().Json(400, http.Json{"error": "Invalid input"})
	}

	if err := facades.Orm().Query().Save(&item); err != nil {
		return ctx.Response().Json(500, http.Json{"error": err.Error()})
	}
	return ctx.Response().Json(200, item)
}

// DELETE /items/{id}
func (c *ItemController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	var item models.Item
	result, err := facades.Orm().Query().Where("id = ?", id).Delete(&item)
	if err != nil {
		return ctx.Response().Json(500, http.Json{
			"error": err.Error(),
		})
	}

	if result.RowsAffected == 0 {
		return ctx.Response().Json(404, http.Json{
			"error": "Item not found",
		})
	}

	return ctx.Response().Json(200, http.Json{
		"message": "Item deleted successfully",
	})
}


// GET /items/search?query=...
func (c *ItemController) Search(ctx http.Context) http.Response {
	query := ctx.Request().Input("query")
	var items []models.Item
	if err := facades.Orm().Query().
		Where("nama_barang LIKE ?", "%"+query+"%").
		OrWhere("tipe_barang LIKE ?", "%"+query+"%").
		Get(&items); err != nil {
		return ctx.Response().Json(500, http.Json{"error": err.Error()})
	}
	return ctx.Response().Json(200, items)
}
