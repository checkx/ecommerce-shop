package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"ecommerce-shop/internal/entity"
	"ecommerce-shop/internal/helpers"
	"ecommerce-shop/internal/service"
)

type WarehousesHandler struct {
	DB  *sqlx.DB
	Svc *service.WarehousesService
}

func (h *WarehousesHandler) Activate(c *gin.Context) {
	id := c.Param("id")
	if err := h.Svc.SetActive(c, id, true); err != nil {
		helpers.WriteError(c.Writer, http.StatusInternalServerError, "DB error", err.Error(), nil)
		return
	}
	helpers.WriteSuccess(c.Writer, "Warehouse activated", nil)
}

func (h *WarehousesHandler) Deactivate(c *gin.Context) {
	id := c.Param("id")
	if err := h.Svc.SetActive(c, id, false); err != nil {
		helpers.WriteError(c.Writer, http.StatusInternalServerError, "DB error", err.Error(), nil)
		return
	}
	helpers.WriteSuccess(c.Writer, "Warehouse deactivated", nil)
}

func (h *WarehousesHandler) Transfer(c *gin.Context) {
	var req entity.TransferReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.WriteError(c.Writer, http.StatusBadRequest, "Invalid JSON", err.Error(), nil)
		return
	}
	err := h.Svc.Transfer(c, req.From, req.To, req.ProductID, req.Quantity)
	if err != nil {
		helpers.WriteError(c.Writer, http.StatusBadRequest, "Transfer failed", err.Error(), nil)
		return
	}
	helpers.WriteSuccess(c.Writer, "Transfer successful", entity.TransferResponse{
		From:      req.From,
		To:        req.To,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Status:    "ok",
	})
}
