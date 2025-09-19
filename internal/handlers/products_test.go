package handlers

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"ecommerce-shop/internal/service"
	"ecommerce-shop/testutils"
)

func TestProductsHandler_ListByShop(t *testing.T) {
	tests := []struct {
		name           string
		shopID         string
		mockSetup      func(sqlmock.Sqlmock)
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "successful list products",
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
			expectedStatus: 200,
		},
		{
			name:   "no products found",
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
			expectedStatus: 200,
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
			expectedStatus: 500,
			expectedError:  "DB error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			db, mock := testutils.MockDB(t)
			defer db.Close()

			productsService := &service.ProductsService{
				DB: db,
			}
			handler := &ProductsHandler{
				DB:  db,
				Svc: productsService,
			}

			tt.mockSetup(mock)

			// Create test context
			c, w := testutils.TestGinContext()
			c.Params = gin.Params{{Key: "shop_id", Value: tt.shopID}}

			// Execute
			handler.ListByShop(c)

			// Assert
			if tt.expectedError != "" {
				testutils.AssertErrorResponse(t, w, tt.expectedStatus, tt.expectedError)
			} else {
				assert.Equal(t, tt.expectedStatus, w.Code)
				assert.Contains(t, w.Body.String(), "Products listed")
			}

			// Verify all expectations
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
