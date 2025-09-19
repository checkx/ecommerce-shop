package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"ecommerce-shop/internal/config"
	"ecommerce-shop/internal/entity"
	"ecommerce-shop/internal/helpers"
	"ecommerce-shop/internal/service"
)

type AuthHandler struct {
	DB       *sqlx.DB
	Log      *zap.Logger
	Validate *validator.Validate
	Cfg      config.Config
	Svc      *service.AuthService
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req entity.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.WriteError(c.Writer, http.StatusBadRequest, "Invalid JSON", err.Error(), h.Log)
		return
	}
	if err := h.Validate.Struct(req); err != nil {
		helpers.WriteError(c.Writer, http.StatusBadRequest, "Validation error", err.Error(), h.Log)
		return
	}
	id, token, err := h.Svc.Register(c, req.Email, req.Password)
	if err != nil {
		helpers.WriteError(c.Writer, http.StatusConflict, "Email exists", err.Error(), h.Log)
		return
	}
	helpers.WriteSuccess(c.Writer, "Register successful", entity.AuthResponse{
		ID:    id,
		Token: token,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req entity.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.WriteError(c.Writer, http.StatusBadRequest, "Invalid JSON", err.Error(), h.Log)
		return
	}
	if err := h.Validate.Struct(req); err != nil {
		helpers.WriteError(c.Writer, http.StatusBadRequest, "Validation error", err.Error(), h.Log)
		return
	}
	id, token, err := h.Svc.Login(c, req.Email, req.Password)
	if err != nil {
		helpers.WriteError(c.Writer, http.StatusUnauthorized, "Invalid credentials", err.Error(), h.Log)
		return
	}
	helpers.WriteSuccess(c.Writer, "Login successful", entity.AuthResponse{
		ID:    id,
		Token: token,
	})
}
