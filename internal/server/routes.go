package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"

	"ecommerce-shop/internal/config"
	"ecommerce-shop/internal/handlers"
	"ecommerce-shop/internal/server/web"
	"ecommerce-shop/internal/service"
)

func RegisterRoutes(r *gin.Engine, cfg config.Config, log *zap.Logger, db *sqlx.DB) {
	api := r.Group("/api")
	{
		v := validator.New()
		authSvc := &service.AuthService{DB: db, Log: log, JWTSecret: cfg.JWTSecret}
		prodSvc := &service.ProductsService{DB: db}
		ordSvc := &service.OrdersService{DB: db, Log: log, TTLMin: cfg.ReservationTTLMinutes}
		whSvc := &service.WarehousesService{DB: db}

		authH := &handlers.AuthHandler{DB: db, Log: log, Validate: v, Cfg: cfg, Svc: authSvc}
		prodH := &handlers.ProductsHandler{DB: db, Svc: prodSvc}
		ordH := &handlers.OrdersHandler{DB: db, Log: log, Validate: v, TTLMin: cfg.ReservationTTLMinutes, Svc: ordSvc}
		whH := &handlers.WarehousesHandler{DB: db, Svc: whSvc}

		// auth
		api.POST("/register", authH.Register)
		api.POST("/login", authH.Login)

		// products
		api.GET("/shops/:shop_id/products", prodH.ListByShop)

		// orders
		api.POST("/orders", ordH.Create)
		api.POST("/orders/:id/pay", ordH.Pay)

		// warehouses
		api.POST("/warehouses/:id/activate", web.JWTAuth(cfg.JWTSecret), whH.Activate)
		api.POST("/warehouses/:id/deactivate", web.JWTAuth(cfg.JWTSecret), whH.Deactivate)
		api.POST("/warehouses/transfer", web.JWTAuth(cfg.JWTSecret), whH.Transfer)
	}
}
