package handlers

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"ecommerce-shop/internal/entity"
	"ecommerce-shop/internal/service"
	"ecommerce-shop/testutils"
)

func TestAuthHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		request        entity.RegisterReq
		mockSetup      func(sqlmock.Sqlmock)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful registration",
			request: entity.RegisterReq{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO users\(email, password_hash\) VALUES \(\$1,\$2\) RETURNING id`).
					WithArgs("test@example.com", sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("user-123"))
			},
			expectedStatus: 200,
		},
		{
			name: "invalid email format",
			request: entity.RegisterReq{
				Email:    "invalid-email",
				Password: "password123",
			},
			mockSetup:      func(mock sqlmock.Sqlmock) {},
			expectedStatus: 400,
			expectedError:  "Validation error",
		},
		{
			name: "password too short",
			request: entity.RegisterReq{
				Email:    "test@example.com",
				Password: "123",
			},
			mockSetup:      func(mock sqlmock.Sqlmock) {},
			expectedStatus: 400,
			expectedError:  "Validation error",
		},
		{
			name: "database error",
			request: entity.RegisterReq{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO users\(email, password_hash\) VALUES \(\$1,\$2\) RETURNING id`).
					WithArgs("test@example.com", sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			expectedStatus: 409,
			expectedError:  "Email exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			db, mock := testutils.MockDB(t)
			defer db.Close()

			logger := testutils.MockLogger(t)
			validator := testutils.TestValidator()
			config := testutils.TestConfig()

			authService := &service.AuthService{
				DB:        db,
				Log:       logger,
				JWTSecret: "test-secret-key",
			}
			handler := &AuthHandler{
				DB:       db,
				Log:      logger,
				Validate: validator,
				Cfg:      config,
				Svc:      authService,
			}

			tt.mockSetup(mock)

			// Create test context
			c, w := testutils.TestGinContextWithBody(t, tt.request)

			// Execute
			handler.Register(c)

			// Assert
			if tt.expectedError != "" {
				testutils.AssertErrorResponse(t, w, tt.expectedStatus, tt.expectedError)
			} else {
				assert.Equal(t, tt.expectedStatus, w.Code)
				assert.Contains(t, w.Body.String(), "Register successful")
			}

			// Verify all expectations
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		request        entity.LoginReq
		mockSetup      func(sqlmock.Sqlmock)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful login",
			request: entity.LoginReq{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Mock password hash for "password123"
				hashedPassword := "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"
				mock.ExpectQuery(`SELECT id, password_hash FROM users WHERE email=\$1`).
					WithArgs("test@example.com").
					WillReturnRows(sqlmock.NewRows([]string{"id", "password_hash"}).AddRow("user-123", hashedPassword))
			},
			expectedStatus: 200,
		},
		{
			name: "invalid email format",
			request: entity.LoginReq{
				Email:    "invalid-email",
				Password: "password123",
			},
			mockSetup:      func(mock sqlmock.Sqlmock) {},
			expectedStatus: 400,
			expectedError:  "Validation error",
		},
		{
			name: "user not found",
			request: entity.LoginReq{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, password_hash FROM users WHERE email=\$1`).
					WithArgs("test@example.com").
					WillReturnError(sql.ErrNoRows)
			},
			expectedStatus: 401,
			expectedError:  "Invalid credentials",
		},
		{
			name: "wrong password",
			request: entity.LoginReq{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Mock password hash for "password123"
				hashedPassword := "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"
				mock.ExpectQuery(`SELECT id, password_hash FROM users WHERE email=\$1`).
					WithArgs("test@example.com").
					WillReturnRows(sqlmock.NewRows([]string{"id", "password_hash"}).AddRow("user-123", hashedPassword))
			},
			expectedStatus: 401,
			expectedError:  "Invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			db, mock := testutils.MockDB(t)
			defer db.Close()

			logger := testutils.MockLogger(t)
			validator := testutils.TestValidator()
			config := testutils.TestConfig()

			authService := &service.AuthService{
				DB:        db,
				Log:       logger,
				JWTSecret: "test-secret-key",
			}
			handler := &AuthHandler{
				DB:       db,
				Log:      logger,
				Validate: validator,
				Cfg:      config,
				Svc:      authService,
			}

			tt.mockSetup(mock)

			// Create test context
			c, w := testutils.TestGinContextWithBody(t, tt.request)

			// Execute
			handler.Login(c)

			// Assert
			if tt.expectedError != "" {
				testutils.AssertErrorResponse(t, w, tt.expectedStatus, tt.expectedError)
			} else {
				assert.Equal(t, tt.expectedStatus, w.Code)
				assert.Contains(t, w.Body.String(), "Login successful")
			}

			// Verify all expectations
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
