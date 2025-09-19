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

func TestWarehousesHandler_Activate(t *testing.T) {
	tests := []struct {
		name           string
		warehouseID    string
		mockSetup      func(sqlmock.Sqlmock)
		expectedStatus int
		expectedError  string
	}{
		{
			name:        "successful activation",
			warehouseID: "wh-123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE warehouses SET active=\$2 WHERE id=\$1`).
					WithArgs("wh-123", true).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedStatus: 200,
		},
		{
			name:        "warehouse not found",
			warehouseID: "wh-123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE warehouses SET active=\$2 WHERE id=\$1`).
					WithArgs("wh-123", true).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedStatus: 200, // Still returns 200 even if no rows affected
		},
		{
			name:        "database error",
			warehouseID: "wh-123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE warehouses SET active=\$2 WHERE id=\$1`).
					WithArgs("wh-123", true).
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

			warehousesService := &service.WarehousesService{
				DB: db,
			}
			handler := &WarehousesHandler{
				DB:  db,
				Svc: warehousesService,
			}

			tt.mockSetup(mock)

			// Create test context
			c, w := testutils.TestGinContext()
			c.Params = gin.Params{{Key: "id", Value: tt.warehouseID}}

			// Execute
			handler.Activate(c)

			// Assert
			if tt.expectedError != "" {
				testutils.AssertErrorResponse(t, w, tt.expectedStatus, tt.expectedError)
			} else {
				assert.Equal(t, tt.expectedStatus, w.Code)
				assert.Contains(t, w.Body.String(), "Warehouse activated")
			}

			// Verify all expectations
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestWarehousesHandler_Deactivate(t *testing.T) {
	tests := []struct {
		name           string
		warehouseID    string
		mockSetup      func(sqlmock.Sqlmock)
		expectedStatus int
		expectedError  string
	}{
		{
			name:        "successful deactivation",
			warehouseID: "wh-123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE warehouses SET active=\$2 WHERE id=\$1`).
					WithArgs("wh-123", false).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedStatus: 200,
		},
		{
			name:        "warehouse not found",
			warehouseID: "wh-123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE warehouses SET active=\$2 WHERE id=\$1`).
					WithArgs("wh-123", false).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedStatus: 200, // Still returns 200 even if no rows affected
		},
		{
			name:        "database error",
			warehouseID: "wh-123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE warehouses SET active=\$2 WHERE id=\$1`).
					WithArgs("wh-123", false).
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

			warehousesService := &service.WarehousesService{
				DB: db,
			}
			handler := &WarehousesHandler{
				DB:  db,
				Svc: warehousesService,
			}

			tt.mockSetup(mock)

			// Create test context
			c, w := testutils.TestGinContext()
			c.Params = gin.Params{{Key: "id", Value: tt.warehouseID}}

			// Execute
			handler.Deactivate(c)

			// Assert
			if tt.expectedError != "" {
				testutils.AssertErrorResponse(t, w, tt.expectedStatus, tt.expectedError)
			} else {
				assert.Equal(t, tt.expectedStatus, w.Code)
				assert.Contains(t, w.Body.String(), "Warehouse deactivated")
			}

			// Verify all expectations
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestWarehousesHandler_Transfer(t *testing.T) {
	tests := []struct {
		name           string
		request        entity.TransferReq
		mockSetup      func(sqlmock.Sqlmock)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful transfer",
			request: entity.TransferReq{
				From:      "wh-1",
				To:        "wh-2",
				ProductID: "prod-1",
				Quantity:  5,
			},
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
			expectedStatus: 200,
		},
		{
			name: "missing required fields",
			request: entity.TransferReq{
				From:      "",
				To:        "wh-2",
				ProductID: "prod-1",
				Quantity:  5,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// No database calls expected for validation error
			},
			expectedStatus: 400,
			expectedError:  "Invalid JSON",
		},
		{
			name: "database error during transfer",
			request: entity.TransferReq{
				From:      "wh-1",
				To:        "wh-2",
				ProductID: "prod-1",
				Quantity:  5,
			},
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
			expectedStatus: 400,
			expectedError:  "Transfer failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			db, mock := testutils.MockDB(t)
			defer db.Close()

			warehousesService := &service.WarehousesService{
				DB: db,
			}
			handler := &WarehousesHandler{
				DB:  db,
				Svc: warehousesService,
			}

			tt.mockSetup(mock)

			// Create test context
			c, w := testutils.TestGinContextWithBody(t, tt.request)

			// Execute
			handler.Transfer(c)

			// Assert
			if tt.expectedError != "" {
				testutils.AssertErrorResponse(t, w, tt.expectedStatus, tt.expectedError)
			} else {
				assert.Equal(t, tt.expectedStatus, w.Code)
				assert.Contains(t, w.Body.String(), "Transfer successful")
			}

			// Verify all expectations
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
