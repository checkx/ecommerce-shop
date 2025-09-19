package service

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type ProductsService struct {
	DB *sqlx.DB
}

type ProductAvailability struct {
	ID, SKU, Name string
	Available     int
}

func (s *ProductsService) ListByShop(ctx context.Context, shopID string) ([]ProductAvailability, error) {
	rows, err := s.DB.QueryxContext(ctx, `
		WITH active_wh as (
			SELECT id FROM warehouses WHERE shop_id=$1 AND active=TRUE
		), inv as (
			SELECT product_id, SUM(quantity) as qty FROM inventory WHERE warehouse_id IN (SELECT id FROM active_wh) GROUP BY product_id
		), res as (
			SELECT product_id, COALESCE(SUM(quantity),0) as reserved
			FROM reservations WHERE released=FALSE AND expires_at>now() AND warehouse_id IN (SELECT id FROM active_wh)
			GROUP BY product_id
		)
		SELECT p.id, p.sku, p.name, COALESCE(inv.qty,0) - COALESCE(res.reserved,0) AS available
		FROM products p
		LEFT JOIN inv ON inv.product_id = p.id
		LEFT JOIN res ON res.product_id = p.id
		ORDER BY p.name
	`, shopID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ProductAvailability
	for rows.Next() {
		var p ProductAvailability
		if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Available); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
