package handlers

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"ecommerce-shop/internal/entity"
	"ecommerce-shop/internal/service"
	"ecommerce-shop/testutils"
)

func TestOrdersHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		idempotencyKey string
		request        entity.CreateOrderReq
		mockSetup      func(sqlmock.Sqlmock)
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "successful order creation",
			idempotencyKey: "test-key-123",
			request: entity.CreateOrderReq{
				ShopID: "shop-123",
				Items: []entity.OrderItemReq{
					{ProductID: "prod-1", Quantity: 2},
					{ProductID: "prod-2", Quantity: 1},
				},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock idempotency key check
				mock.ExpectQuery(`SELECT order_id FROM idempotency_keys WHERE key=\$1 AND request_hash=\$2`).
					WithArgs("test-key-123", sqlmock.AnyArg()).
					WillReturnError(sql.ErrNoRows)

				// Mock idempotency key insert
				mock.ExpectExec(`INSERT INTO idempotency_keys\(key, user_id, request_hash\) VALUES \(\$1, gen_random_uuid\(\), \$2\)`).
					WithArgs("test-key-123", sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Mock order creation
				mock.ExpectQuery(`INSERT INTO orders\(user_id, shop_id, status\) VALUES \(gen_random_uuid\(\), \$1, 'reserved'\) RETURNING id`).
					WithArgs("shop-123").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("order-123"))

				// Mock order items insert
				mock.ExpectExec(`INSERT INTO order_items\(order_id, product_id, quantity\) VALUES \(\$1,\$2,\$3\)`).
					WithArgs("order-123", "prod-1", 2).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Mock active warehouses query
				mock.ExpectQuery(`SELECT id FROM warehouses WHERE shop_id=\$1 AND active=TRUE`).
					WithArgs("shop-123").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("wh-1"))

				// Mock inventory lock
				mock.ExpectQuery(`SELECT quantity FROM inventory WHERE warehouse_id=\$1 AND product_id=\$2 FOR UPDATE`).
					WithArgs("wh-1", "prod-1").
					WillReturnRows(sqlmock.NewRows([]string{"quantity"}).AddRow(10))

				// Mock reserved quantity query
				mock.ExpectQuery(`SELECT COALESCE\(SUM\(quantity\),0\) FROM reservations WHERE warehouse_id=\$1 AND product_id=\$2 AND released=FALSE AND expires_at>now\(\)`).
					WithArgs("wh-1", "prod-1").
					WillReturnRows(sqlmock.NewRows([]string{"sum"}).AddRow(0))

				// Mock reservation insert
				mock.ExpectExec(`INSERT INTO reservations\(order_id, warehouse_id, product_id, quantity, expires_at\) VALUES \(\$1,\$2,\$3,\$4,now\(\)\+interval '\$5 minutes'\)`).
					WithArgs("order-123", "wh-1", "prod-1", 2, 15).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Mock second order item
				mock.ExpectExec(`INSERT INTO order_items\(order_id, product_id, quantity\) VALUES \(\$1,\$2,\$3\)`).
					WithArgs("order-123", "prod-2", 1).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Mock active warehouses query for second item
				mock.ExpectQuery(`SELECT id FROM warehouses WHERE shop_id=\$1 AND active=TRUE`).
					WithArgs("shop-123").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("wh-1"))

				// Mock inventory lock for second item
				mock.ExpectQuery(`SELECT quantity FROM inventory WHERE warehouse_id=\$1 AND product_id=\$2 FOR UPDATE`).
					WithArgs("wh-1", "prod-2").
					WillReturnRows(sqlmock.NewRows([]string{"quantity"}).AddRow(5))

				// Mock reserved quantity query for second item
				mock.ExpectQuery(`SELECT COALESCE\(SUM\(quantity\),0\) FROM reservations WHERE warehouse_id=\$1 AND product_id=\$2 AND released=FALSE AND expires_at>now\(\)`).
					WithArgs("wh-1", "prod-2").
					WillReturnRows(sqlmock.NewRows([]string{"sum"}).AddRow(0))

				// Mock reservation insert for second item
				mock.ExpectExec(`INSERT INTO reservations\(order_id, warehouse_id, product_id, quantity, expires_at\) VALUES \(\$1,\$2,\$3,\$4,now\(\)\+interval '\$5 minutes'\)`).
					WithArgs("order-123", "wh-1", "prod-2", 1, 15).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Mock idempotency key update
				mock.ExpectExec(`UPDATE idempotency_keys SET order_id=\$2 WHERE key=\$1`).
					WithArgs("test-key-123", "order-123").
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Mock transaction commit
				mock.ExpectCommit()
			},
			expectedStatus: 200,
		},
		{
			name:           "missing idempotency key",
			idempotencyKey: "",
			request: entity.CreateOrderReq{
				ShopID: "shop-123",
				Items:  []entity.OrderItemReq{{ProductID: "prod-1", Quantity: 1}},
			},
			mockSetup:      func(mock sqlmock.Sqlmock) {},
			expectedStatus: 400,
			expectedError:  "Missing Idempotency-Key",
		},
		{
			name:           "invalid JSON",
			idempotencyKey: "test-key-123",
			request:        entity.CreateOrderReq{},
			mockSetup:      func(mock sqlmock.Sqlmock) {},
			expectedStatus: 400,
			expectedError:  "Invalid JSON",
		},
		{
			name:           "validation error - missing shop_id",
			idempotencyKey: "test-key-123",
			request: entity.CreateOrderReq{
				Items: []entity.OrderItemReq{{ProductID: "prod-1", Quantity: 1}},
			},
			mockSetup:      func(mock sqlmock.Sqlmock) {},
			expectedStatus: 400,
			expectedError:  "Validation error",
		},
		{
			name:           "validation error - empty items",
			idempotencyKey: "test-key-123",
			request: entity.CreateOrderReq{
				ShopID: "shop-123",
				Items:  []entity.OrderItemReq{},
			},
			mockSetup:      func(mock sqlmock.Sqlmock) {},
			expectedStatus: 400,
			expectedError:  "Validation error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			db, mock := testutils.MockDB(t)
			defer db.Close()

			logger := testutils.MockLogger(t)
			validator := testutils.TestValidator()

			ordersService := &service.OrdersService{
				DB:     db,
				Log:    logger,
				TTLMin: 15,
			}
			handler := &OrdersHandler{
				DB:       db,
				Log:      logger,
				Validate: validator,
				TTLMin:   15,
				Svc:      ordersService,
			}

			tt.mockSetup(mock)

			// Create test context
			c, w := testutils.TestGinContextWithBody(t, tt.request)
			if tt.idempotencyKey != "" {
				c.Request.Header.Set("Idempotency-Key", tt.idempotencyKey)
			}

			// Execute
			handler.Create(c)

			// Assert
			if tt.expectedError != "" {
				testutils.AssertErrorResponse(t, w, tt.expectedStatus, tt.expectedError)
			} else {
				assert.Equal(t, tt.expectedStatus, w.Code)
				assert.Contains(t, w.Body.String(), "Order reserved")
			}

			// Verify all expectations
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOrdersHandler_Pay(t *testing.T) {
	tests := []struct {
		name           string
		orderID        string
		mockSetup      func(sqlmock.Sqlmock)
		expectedStatus int
		expectedError  string
	}{
		{
			name:    "successful payment",
			orderID: "order-123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock order status check
				mock.ExpectQuery(`SELECT status FROM orders WHERE id=\$1 FOR UPDATE`).
					WithArgs("order-123").
					WillReturnRows(sqlmock.NewRows([]string{"status"}).AddRow("reserved"))

				// Mock reservations query
				mock.ExpectQuery(`SELECT warehouse_id, product_id, quantity FROM reservations WHERE order_id=\$1 AND released=FALSE AND expires_at>now\(\)`).
					WithArgs("order-123").
					WillReturnRows(sqlmock.NewRows([]string{"warehouse_id", "product_id", "quantity"}).
						AddRow("wh-1", "prod-1", 2).
						AddRow("wh-1", "prod-2", 1))

				// Mock inventory updates
				mock.ExpectExec(`UPDATE inventory SET quantity = quantity - \$3 WHERE warehouse_id=\$1 AND product_id=\$2`).
					WithArgs("wh-1", "prod-1", 2).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec(`UPDATE inventory SET quantity = quantity - \$3 WHERE warehouse_id=\$1 AND product_id=\$2`).
					WithArgs("wh-1", "prod-2", 1).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Mock reservations release
				mock.ExpectExec(`UPDATE reservations SET released=TRUE WHERE order_id=\$1`).
					WithArgs("order-123").
					WillReturnResult(sqlmock.NewResult(1, 2))

				// Mock order status update
				mock.ExpectExec(`UPDATE orders SET status='paid', updated_at=now\(\) WHERE id=\$1`).
					WithArgs("order-123").
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Mock transaction commit
				mock.ExpectCommit()
			},
			expectedStatus: 200,
		},
		{
			name:    "order not found",
			orderID: "order-123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock order status check - order not found
				mock.ExpectQuery(`SELECT status FROM orders WHERE id=\$1 FOR UPDATE`).
					WithArgs("order-123").
					WillReturnError(sql.ErrNoRows)

				// Mock transaction rollback
				mock.ExpectRollback()
			},
			expectedStatus: 400,
			expectedError:  "Cannot pay",
		},
		{
			name:    "order not in reserved status",
			orderID: "order-123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock order status check - order already paid
				mock.ExpectQuery(`SELECT status FROM orders WHERE id=\$1 FOR UPDATE`).
					WithArgs("order-123").
					WillReturnRows(sqlmock.NewRows([]string{"status"}).AddRow("paid"))

				// Mock transaction rollback
				mock.ExpectRollback()
			},
			expectedStatus: 400,
			expectedError:  "Cannot pay",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			db, mock := testutils.MockDB(t)
			defer db.Close()

			logger := testutils.MockLogger(t)

			ordersService := &service.OrdersService{
				DB:     db,
				Log:    logger,
				TTLMin: 15,
			}
			handler := &OrdersHandler{
				DB:       db,
				Log:      logger,
				Validate: testutils.TestValidator(),
				TTLMin:   15,
				Svc:      ordersService,
			}

			tt.mockSetup(mock)

			// Create test context
			c, w := testutils.TestGinContext()
			c.Params = gin.Params{{Key: "id", Value: tt.orderID}}

			// Execute
			handler.Pay(c)

			// Assert
			if tt.expectedError != "" {
				testutils.AssertErrorResponse(t, w, tt.expectedStatus, tt.expectedError)
			} else {
				assert.Equal(t, tt.expectedStatus, w.Code)
				assert.Contains(t, w.Body.String(), "Order paid")
			}

			// Verify all expectations
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
