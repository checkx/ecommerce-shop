package db

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"ecommerce-shop/internal/config"
)

func Connect(cfg config.Config, log *zap.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", cfg.DBURL)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.Info("connected to postgres")
	return db, nil
}
