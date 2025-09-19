package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"ecommerce-shop/internal/config"
	"ecommerce-shop/internal/db"
	"ecommerce-shop/internal/logger"
	"ecommerce-shop/internal/server"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg)

	database, err := db.Connect(cfg, log)
	if err != nil {
		log.Fatal("failed to connect to db", zap.Error(err))
	}
	defer func(dbx *sqlx.DB) { _ = dbx.Close() }(database)

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := server.BuildRouter(cfg, log, database)

	srv := server.NewHTTPServer(cfg, log, r)

	// graceful shutdown
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatal("server error", zap.Error(err))
		}
	}()

	// start background worker
	// bgCtx, bgCancel := context.WithCancel(context.Background())
	// defer bgCancel()
	// go worker.NewReleaser(database, log, 30*time.Second).Start(bgCtx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("graceful shutdown failed", zap.Error(err))
	}
}
