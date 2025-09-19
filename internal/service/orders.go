package service

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"ecommerce-shop/internal/repo"
)

type OrdersService struct {
	DB     *sqlx.DB
	Log    *zap.Logger
	TTLMin int
}

func hash(b []byte) string { h := sha256.Sum256(b); return hex.EncodeToString(h[:]) }

func (s *OrdersService) Create(ctx context.Context, idempotencyKey string, rawBody []byte, shopID string, items []struct {
	ProductID string
	Quantity  int
}) (string, error) {
	var orderID string
	requestHash := hash(rawBody)
	err := repo.New(s.DB).WithTx(ctx, func(tx *sqlx.Tx) error {
		var existing sql.NullString
		if err := tx.QueryRowxContext(ctx, `SELECT order_id FROM idempotency_keys WHERE key=$1 AND request_hash=$2`, idempotencyKey, requestHash).Scan(&existing); err == nil {
			if existing.Valid {
				orderID = existing.String
				return nil
			}
		} else if err != sql.ErrNoRows {
			return err
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO idempotency_keys(key, user_id, request_hash) VALUES ($1, gen_random_uuid(), $2)`, idempotencyKey, requestHash); err != nil {
			return err
		}
		if err := tx.GetContext(ctx, &orderID, `INSERT INTO orders(user_id, shop_id, status) VALUES (gen_random_uuid(), $1, 'reserved') RETURNING id`, shopID); err != nil {
			return err
		}
		for _, it := range items {
			if _, err := tx.ExecContext(ctx, `INSERT INTO order_items(order_id, product_id, quantity) VALUES ($1,$2,$3)`, orderID, it.ProductID, it.Quantity); err != nil {
				return err
			}
			whIDs, err := repo.ActiveWarehousesForShop(ctx, tx, shopID)
			if err != nil {
				return err
			}
			reserved := false
			for _, wh := range whIDs {
				invQty, err := repo.LockInventoryRow(ctx, tx, wh, it.ProductID)
				if err != nil {
					return err
				}
				resQty, err := repo.SumReservedNotExpired(ctx, tx, wh, it.ProductID)
				if err != nil {
					return err
				}
				if invQty-resQty >= it.Quantity {
					if err := repo.ReserveStock(ctx, tx, orderID, wh, it.ProductID, it.Quantity, s.TTLMin); err != nil {
						return err
					}
					reserved = true
					break
				}
			}
			if !reserved {
				return sql.ErrNoRows
			}
		}
		_, err := tx.ExecContext(ctx, `UPDATE idempotency_keys SET order_id=$2 WHERE key=$1`, idempotencyKey, orderID)
		return err
	})
	return orderID, err
}

func (s *OrdersService) Pay(ctx context.Context, orderID string) error {
	return repo.New(s.DB).WithTx(ctx, func(tx *sqlx.Tx) error {
		var status string
		if err := tx.GetContext(ctx, &status, `SELECT status FROM orders WHERE id=$1 FOR UPDATE`, orderID); err != nil {
			return err
		}
		if status != "reserved" {
			return sql.ErrNoRows
		}
		rows, err := tx.QueryxContext(ctx, `SELECT warehouse_id, product_id, quantity FROM reservations WHERE order_id=$1 AND released=FALSE AND expires_at>now()`, orderID)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var wh, pid string
			var qty int
			if err := rows.Scan(&wh, &pid, &qty); err != nil {
				return err
			}
			if _, err := tx.ExecContext(ctx, `UPDATE inventory SET quantity = quantity - $3 WHERE warehouse_id=$1 AND product_id=$2`, wh, pid, qty); err != nil {
				return err
			}
		}
		if _, err := tx.ExecContext(ctx, `UPDATE reservations SET released=TRUE WHERE order_id=$1`, orderID); err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, `UPDATE orders SET status='paid', updated_at=now() WHERE id=$1`, orderID)
		return err
	})
}
