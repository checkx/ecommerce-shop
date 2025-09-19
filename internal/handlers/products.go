package handlers

import (
	"net/http"

	"ecommerce-shop/internal/entity"
	"ecommerce-shop/internal/helpers"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"ecommerce-shop/internal/service"
)

type ProductsHandler struct {
	DB  *sqlx.DB
	Svc *service.ProductsService
}

func (h *ProductsHandler) ListByShop(c *gin.Context) {
	shopID := c.Param("shop_id")
	list, err := h.Svc.ListByShop(c, shopID)
	if err != nil {
		helpers.WriteError(c.Writer, http.StatusInternalServerError, "DB error", err.Error(), nil)
		return
	}
	out := make([]entity.ProductResponse, 0, len(list))
	for _, p := range list {
		out = append(out, entity.ProductResponse{
			ID:        p.ID,
			SKU:       p.SKU,
			Name:      p.Name,
			Available: p.Available,
		})
	}
	helpers.WriteSuccess(c.Writer, "Products listed", out)
}
