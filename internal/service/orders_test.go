package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"ecommerce-shop/testutils"
)

func TestOrdersService_Pay(t *testing.T) {
	tests := []struct {
		name      string
		orderID   string
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
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
			wantErr: false,
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
			wantErr: true,
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
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			db, mock := testutils.MockDB(t)
			defer db.Close()

			logger := testutils.MockLogger(t)
			service := &OrdersService{
				DB:     db,
				Log:    logger,
				TTLMin: 15,
			}

			tt.mockSetup(mock)

			// Execute
			err := service.Pay(context.Background(), tt.orderID)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Verify all expectations
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
