package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"ecommerce-shop/testutils"
)

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		password  string
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name:     "successful registration",
			email:    "test@example.com",
			password: "password123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO users\(email, password_hash\) VALUES \(\$1,\$2\) RETURNING id`).
					WithArgs("test@example.com", sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("user-123"))
			},
			wantErr: false,
		},
		{
			name:     "duplicate email",
			email:    "test@example.com",
			password: "password123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO users\(email, password_hash\) VALUES \(\$1,\$2\) RETURNING id`).
					WithArgs("test@example.com", sqlmock.AnyArg()).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
		{
			name:     "database error",
			email:    "test@example.com",
			password: "password123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO users\(email, password_hash\) VALUES \(\$1,\$2\) RETURNING id`).
					WithArgs("test@example.com", sqlmock.AnyArg()).
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

			logger := testutils.MockLogger(t)
			service := &AuthService{
				DB:        db,
				Log:       logger,
				JWTSecret: "test-secret-key",
			}

			tt.mockSetup(mock)

			// Execute
			id, token, err := service.Register(context.Background(), tt.email, tt.password)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, id)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, id)
				assert.NotEmpty(t, token)
			}

			// Verify all expectations
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
