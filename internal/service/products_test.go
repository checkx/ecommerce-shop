package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"ecommerce-shop/testutils"
)

func TestProductsService_ListByShop(t *testing.T) {
	tests := []struct {
		name      string
		shopID    string
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
		wantLen   int
	}{
		{
			name:   "successful list with products",
			shopID: "shop-123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "sku", "name", "available"}).
					AddRow("prod-1", "SKU001", "Product 1", 10).
					AddRow("prod-2", "SKU002", "Product 2", 5)

				mock.ExpectQuery(`WITH active_wh as \(
			SELECT id FROM warehouses WHERE shop_id=\$1 AND active=TRUE
		\), inv as \(
			SELECT product_id, SUM\(quantity\) as qty FROM inventory WHERE warehouse_id IN \(SELECT id FROM active_wh\) GROUP BY product_id
		\), res as \(
			SELECT product_id, COALESCE\(SUM\(quantity\),0\) as reserved
			FROM reservations WHERE released=FALSE AND expires_at>now\(\) AND warehouse_id IN \(SELECT id FROM active_wh\)
			GROUP BY product_id
		\)
		SELECT p\.id, p\.sku, p\.name, COALESCE\(inv\.qty,0\) - COALESCE\(res\.reserved,0\) AS available
		FROM products p
		LEFT JOIN inv ON inv\.product_id = p\.id
		LEFT JOIN res ON res\.product_id = p\.id
		ORDER BY p\.name`).
					WithArgs("shop-123").
					WillReturnRows(rows)
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name:   "successful list with no products",
			shopID: "shop-123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "sku", "name", "available"})

				mock.ExpectQuery(`WITH active_wh as \(
			SELECT id FROM warehouses WHERE shop_id=\$1 AND active=TRUE
		\), inv as \(
			SELECT product_id, SUM\(quantity\) as qty FROM inventory WHERE warehouse_id IN \(SELECT id FROM active_wh\) GROUP BY product_id
		\), res as \(
			SELECT product_id, COALESCE\(SUM\(quantity\),0\) as reserved
			FROM reservations WHERE released=FALSE AND expires_at>now\(\) AND warehouse_id IN \(SELECT id FROM active_wh\)
			GROUP BY product_id
		\)
		SELECT p\.id, p\.sku, p\.name, COALESCE\(inv\.qty,0\) - COALESCE\(res\.reserved,0\) AS available
		FROM products p
		LEFT JOIN inv ON inv\.product_id = p\.id
		LEFT JOIN res ON res\.product_id = p\.id
		ORDER BY p\.name`).
					WithArgs("shop-123").
					WillReturnRows(rows)
			},
			wantErr: false,
			wantLen: 0,
		},
		{
			name:   "database error",
			shopID: "shop-123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`WITH active_wh as \(
			SELECT id FROM warehouses WHERE shop_id=\$1 AND active=TRUE
		\), inv as \(
			SELECT product_id, SUM\(quantity\) as qty FROM inventory WHERE warehouse_id IN \(SELECT id FROM active_wh\) GROUP BY product_id
		\), res as \(
			SELECT product_id, COALESCE\(SUM\(quantity\),0\) as reserved
			FROM reservations WHERE released=FALSE AND expires_at>now\(\) AND warehouse_id IN \(SELECT id FROM active_wh\)
			GROUP BY product_id
		\)
		SELECT p\.id, p\.sku, p\.name, COALESCE\(inv\.qty,0\) - COALESCE\(res\.reserved,0\) AS available
		FROM products p
		LEFT JOIN inv ON inv\.product_id = p\.id
		LEFT JOIN res ON res\.product_id = p\.id
		ORDER BY p\.name`).
					WithArgs("shop-123").
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
			wantLen: 0,
		},
		{
			name:   "scan error",
			shopID: "shop-123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Return rows with wrong number of columns to cause scan error
				rows := sqlmock.NewRows([]string{"id", "sku", "name"}).
					AddRow("prod-1", "SKU001", "Product 1")

				mock.ExpectQuery(`WITH active_wh as \(
			SELECT id FROM warehouses WHERE shop_id=\$1 AND active=TRUE
		\), inv as \(
			SELECT product_id, SUM\(quantity\) as qty FROM inventory WHERE warehouse_id IN \(SELECT id FROM active_wh\) GROUP BY product_id
		\), res as \(
			SELECT product_id, COALESCE\(SUM\(quantity\),0\) as reserved
			FROM reservations WHERE released=FALSE AND expires_at>now\(\) AND warehouse_id IN \(SELECT id FROM active_wh\)
			GROUP BY product_id
		\)
		SELECT p\.id, p\.sku, p\.name, COALESCE\(inv\.qty,0\) - COALESCE\(res\.reserved,0\) AS available
		FROM products p
		LEFT JOIN inv ON inv\.product_id = p\.id
		LEFT JOIN res ON res\.product_id = p\.id
		ORDER BY p\.name`).
					WithArgs("shop-123").
					WillReturnRows(rows)
			},
			wantErr: true,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			db, mock := testutils.MockDB(t)
			defer db.Close()

			service := &ProductsService{
				DB: db,
			}

			tt.mockSetup(mock)

			// Execute
			products, err := service.ListByShop(context.Background(), tt.shopID)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, products, tt.wantLen)
			}

			// Verify all expectations
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
