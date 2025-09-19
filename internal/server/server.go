package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"ecommerce-shop/internal/config"
	"ecommerce-shop/internal/server/web"
)

type HTTPServer struct {
	server *http.Server
	log    *zap.Logger
}

func NewHTTPServer(cfg config.Config, log *zap.Logger, r *gin.Engine) *HTTPServer {
	return &HTTPServer{server: &http.Server{Addr: cfg.HTTPAddr, Handler: r}, log: log}
}

func (s *HTTPServer) Start() error {
	s.log.Info("http server starting", zap.String("addr", s.server.Addr))
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	s.log.Info("http server shutting down")
	return s.server.Shutdown(ctx)
}

func BuildRouter(cfg config.Config, log *zap.Logger, db *sqlx.DB) *gin.Engine {
	r := gin.New()
	r.Use(web.RequestID())
	r.Use(web.ZapLogger(log))
	r.Use(gin.Recovery())

	api := r.Group("/api")
	{
		api.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now().UTC()}) })
	}

	RegisterRoutes(r, cfg, log, db)
	return r
}
