package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"ecommerce-shop/testutils"
)

func TestWarehousesService_SetActive(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		active    bool
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name:   "successful activation",
			id:     "wh-123",
			active: true,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE warehouses SET active=\$2 WHERE id=\$1`).
					WithArgs("wh-123", true).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name:   "successful deactivation",
			id:     "wh-123",
			active: false,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE warehouses SET active=\$2 WHERE id=\$1`).
					WithArgs("wh-123", false).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name:   "warehouse not found",
			id:     "wh-123",
			active: true,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE warehouses SET active=\$2 WHERE id=\$1`).
					WithArgs("wh-123", true).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: false, // No error even if no rows affected
		},
		{
			name:   "database error",
			id:     "wh-123",
			active: true,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE warehouses SET active=\$2 WHERE id=\$1`).
					WithArgs("wh-123", true).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			db, mock := testutils.MockDB(t)
			defer db.Close()

			service := &WarehousesService{
				DB: db,
			}

			tt.mockSetup(mock)

			// Execute
			err := service.SetActive(context.Background(), tt.id, tt.active)

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

func TestWarehousesService_Transfer(t *testing.T) {
	tests := []struct {
		name      string
		from      string
		to        string
		productID string
		qty       int
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name:      "successful transfer",
			from:      "wh-1",
			to:        "wh-2",
			productID: "prod-1",
			qty:       5,
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock inventory reduction
				mock.ExpectExec(`UPDATE inventory SET quantity = quantity - \$3 WHERE warehouse_id=\$1 AND product_id=\$2`).
					WithArgs("wh-1", "prod-1", 5).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Mock inventory addition
				mock.ExpectExec(`INSERT INTO inventory\(warehouse_id, product_id, quantity\) VALUES \(\$1,\$2,\$3\)
		ON CONFLICT \(warehouse_id, product_id\) DO UPDATE SET quantity = inventory\.quantity \+ EXCLUDED\.quantity`).
					WithArgs("wh-2", "prod-1", 5).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Mock transaction commit
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:      "transfer to same warehouse",
			from:      "wh-1",
			to:        "wh-1",
			productID: "prod-1",
			qty:       5,
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock inventory reduction
				mock.ExpectExec(`UPDATE inventory SET quantity = quantity - \$3 WHERE warehouse_id=\$1 AND product_id=\$2`).
					WithArgs("wh-1", "prod-1", 5).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Mock inventory addition (same warehouse)
				mock.ExpectExec(`INSERT INTO inventory\(warehouse_id, product_id, quantity\) VALUES \(\$1,\$2,\$3\)
		ON CONFLICT \(warehouse_id, product_id\) DO UPDATE SET quantity = inventory\.quantity \+ EXCLUDED\.quantity`).
					WithArgs("wh-1", "prod-1", 5).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Mock transaction commit
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:      "database error during reduction",
			from:      "wh-1",
			to:        "wh-2",
			productID: "prod-1",
			qty:       5,
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock inventory reduction - database error
				mock.ExpectExec(`UPDATE inventory SET quantity = quantity - \$3 WHERE warehouse_id=\$1 AND product_id=\$2`).
					WithArgs("wh-1", "prod-1", 5).
					WillReturnError(sql.ErrConnDone)

				// Mock transaction rollback
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name:      "database error during addition",
			from:      "wh-1",
			to:        "wh-2",
			productID: "prod-1",
			qty:       5,
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock inventory reduction
				mock.ExpectExec(`UPDATE inventory SET quantity = quantity - \$3 WHERE warehouse_id=\$1 AND product_id=\$2`).
					WithArgs("wh-1", "prod-1", 5).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Mock inventory addition - database error
				mock.ExpectExec(`INSERT INTO inventory\(warehouse_id, product_id, quantity\) VALUES \(\$1,\$2,\$3\)
		ON CONFLICT \(warehouse_id, product_id\) DO UPDATE SET quantity = inventory\.quantity \+ EXCLUDED\.quantity`).
					WithArgs("wh-2", "prod-1", 5).
					WillReturnError(sql.ErrConnDone)

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

			service := &WarehousesService{
				DB: db,
			}

			tt.mockSetup(mock)

			// Execute
			err := service.Transfer(context.Background(), tt.from, tt.to, tt.productID, tt.qty)

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
