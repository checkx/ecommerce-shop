package service

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type WarehousesService struct{ DB *sqlx.DB }

func (s *WarehousesService) SetActive(ctx context.Context, id string, active bool) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE warehouses SET active=$2 WHERE id=$1`, id, active)
	return err
}

func (s *WarehousesService) Transfer(ctx context.Context, from, to, productID string, qty int) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `UPDATE inventory SET quantity = quantity - $3 WHERE warehouse_id=$1 AND product_id=$2`, from, productID, qty); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO inventory(warehouse_id, product_id, quantity) VALUES ($1,$2,$3)
		ON CONFLICT (warehouse_id, product_id) DO UPDATE SET quantity = inventory.quantity + EXCLUDED.quantity`, to, productID, qty); err != nil {
		return err
	}
	return tx.Commit()
}

