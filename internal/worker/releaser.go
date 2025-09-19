package worker

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"ecommerce-shop/internal/repo"
)

type Releaser struct {
	DB     *sqlx.DB
	Log    *zap.Logger
	Ticker *time.Ticker
}

func NewReleaser(db *sqlx.DB, log *zap.Logger, interval time.Duration) *Releaser {
	return &Releaser{DB: db, Log: log, Ticker: time.NewTicker(interval)}
}

func (w *Releaser) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			w.Ticker.Stop()
			return
		case <-w.Ticker.C:
			count, err := repo.ReleaseExpiredReservations(ctx, w.DB, 500)
			if err != nil {
				w.Log.Error("release reservations failed", zap.Error(err))
				continue
			}
			if count > 0 {
				w.Log.Info("released expired reservations", zap.Int("count", count))
			}
		}
	}
}
