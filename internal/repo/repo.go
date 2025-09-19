package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	DB *sqlx.DB
}

func New(db *sqlx.DB) *Repositories {
	return &Repositories{DB: db}
}

func (r *Repositories) WithTx(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := r.DB.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	if err = fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func LockInventoryRow(ctx context.Context, tx *sqlx.Tx, warehouseID, productID string) (int, error) {
	var qty int

	if err := tx.GetContext(ctx, &qty, `
		WITH up AS (
			INSERT INTO inventory(warehouse_id, product_id, quantity)
			VALUES ($1, $2, 0)
			ON CONFLICT (warehouse_id, product_id) DO NOTHING
			RETURNING quantity
		)
		SELECT quantity FROM inventory WHERE warehouse_id = $1 AND product_id = $2 FOR UPDATE
	`, warehouseID, productID); err != nil {
		return 0, err
	}
	return qty, nil
}

func ActiveWarehousesForShop(ctx context.Context, q sqlx.ExtContext, shopID string) ([]string, error) {
	rows, err := q.QueryxContext(ctx, `SELECT id FROM warehouses WHERE shop_id=$1 AND active=TRUE`, shopID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func ReserveStock(ctx context.Context, tx *sqlx.Tx, orderID, warehouseID, productID string, qty int, ttlMinutes int) error {

	expires := time.Now().Add(time.Duration(ttlMinutes) * time.Minute)
	_, err := tx.ExecContext(ctx, `
		INSERT INTO reservations(order_id, warehouse_id, product_id, quantity, expires_at)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (order_id, warehouse_id, product_id) DO UPDATE SET quantity=EXCLUDED.quantity, expires_at=EXCLUDED.expires_at, released=FALSE
	`, orderID, warehouseID, productID, qty, expires)
	return err
}

func SumReservedNotExpired(ctx context.Context, q sqlx.ExtContext, warehouseID, productID string) (int, error) {
	var sum sql.NullInt64
	if err := sqlx.GetContext(ctx, q, &sum, `
		SELECT COALESCE(SUM(quantity),0) FROM reservations
		WHERE warehouse_id=$1 AND product_id=$2 AND released=FALSE AND expires_at>now()
	`, warehouseID, productID); err != nil {
		return 0, err
	}
	if !sum.Valid {
		return 0, nil
	}
	return int(sum.Int64), nil
}

func ReleaseExpiredReservations(ctx context.Context, db *sqlx.DB, limit int) (int, error) {
	res, err := db.ExecContext(ctx, `UPDATE reservations SET released=TRUE WHERE released=FALSE AND expires_at<=now() LIMIT $1`, limit)
	if err != nil {
		return 0, err
	}
	affected, _ := res.RowsAffected()
	return int(affected), nil
}
