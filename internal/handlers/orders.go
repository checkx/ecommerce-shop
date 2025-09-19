package handlers

import (
	"net/http"

	"ecommerce-shop/internal/helpers"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"ecommerce-shop/internal/entity"
	"ecommerce-shop/internal/service"
)

type OrdersHandler struct {
	DB       *sqlx.DB
	Log      *zap.Logger
	Validate *validator.Validate
	TTLMin   int
	Svc      *service.OrdersService
}

func (h *OrdersHandler) Create(c *gin.Context) {
	idk := c.GetHeader("Idempotency-Key")
	if idk == "" {
		helpers.WriteError(c.Writer, http.StatusBadRequest, "Missing Idempotency-Key", "", h.Log)
		return
	}
	body, err := c.GetRawData()
	if err != nil {
		helpers.WriteError(c.Writer, http.StatusBadRequest, "Bad body", err.Error(), h.Log)
		return
	}
	c.Request.Body = http.NoBody
	var req entity.CreateOrderReq
	if err := c.BindJSON(&req); err != nil {
		helpers.WriteError(c.Writer, http.StatusBadRequest, "Invalid JSON", err.Error(), h.Log)
		return
	}
	if err := h.Validate.Struct(req); err != nil {
		helpers.WriteError(c.Writer, http.StatusBadRequest, "Validation error", err.Error(), h.Log)
		return
	}
	orderID, err := h.Svc.Create(c, idk, body, req.ShopID, func() []struct {
		ProductID string
		Quantity  int
	} {
		out := make([]struct {
			ProductID string
			Quantity  int
		}, 0, len(req.Items))
		for _, it := range req.Items {
			out = append(out, struct {
				ProductID string
				Quantity  int
			}{ProductID: it.ProductID, Quantity: it.Quantity})
		}
		return out
	}())
	if err != nil {
		helpers.WriteError(c.Writer, http.StatusConflict, "Insufficient stock", err.Error(), h.Log)
		return
	}
	helpers.WriteSuccess(c.Writer, "Order reserved", entity.OrderResponse{
		ID:     orderID,
		Status: "reserved",
	})
}

func (h *OrdersHandler) Pay(c *gin.Context) {
	orderID := c.Param("id")
	err := h.Svc.Pay(c, orderID)
	if err != nil {
		helpers.WriteError(c.Writer, http.StatusBadRequest, "Cannot pay", err.Error(), h.Log)
		return
	}
	helpers.WriteSuccess(c.Writer, "Order paid", entity.OrderResponse{
		ID:     orderID,
		Status: "paid",
	})
}
