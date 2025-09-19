package helpers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	"ecommerce-shop/internal/entity"
)

func TestWriteSuccess(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		data     interface{}
		wantCode int
	}{
		{
			name:     "success with data",
			message:  "Operation successful",
			data:     entity.AuthResponse{ID: "user-123", Token: "jwt-token"},
			wantCode: http.StatusOK,
		},
		{
			name:     "success without data",
			message:  "Operation successful",
			data:     nil,
			wantCode: http.StatusOK,
		},
		{
			name:     "success with string data",
			message:  "Operation successful",
			data:     "simple string",
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			w := httptest.NewRecorder()

			// Execute
			WriteSuccess(w, tt.message, tt.data)

			// Assert
			assert.Equal(t, tt.wantCode, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Contains(t, w.Body.String(), "success")
			assert.Contains(t, w.Body.String(), tt.message)
		})
	}
}

func TestWriteError(t *testing.T) {
	tests := []struct {
		name      string
		code      int
		message   string
		errDetail string
		wantCode  int
	}{
		{
			name:      "bad request error",
			code:      http.StatusBadRequest,
			message:   "Validation error",
			errDetail: "invalid input",
			wantCode:  http.StatusBadRequest,
		},
		{
			name:      "internal server error",
			code:      http.StatusInternalServerError,
			message:   "Database error",
			errDetail: "connection failed",
			wantCode:  http.StatusInternalServerError,
		},
		{
			name:      "unauthorized error",
			code:      http.StatusUnauthorized,
			message:   "Invalid credentials",
			errDetail: "wrong password",
			wantCode:  http.StatusUnauthorized,
		},
		{
			name:      "error without detail",
			code:      http.StatusNotFound,
			message:   "Not found",
			errDetail: "",
			wantCode:  http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			w := httptest.NewRecorder()
			logger := zaptest.NewLogger(t)

			// Execute
			WriteError(w, tt.code, tt.message, tt.errDetail, logger)

			// Assert
			assert.Equal(t, tt.wantCode, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Contains(t, w.Body.String(), "error")
			assert.Contains(t, w.Body.String(), tt.message)
			if tt.errDetail != "" {
				assert.Contains(t, w.Body.String(), tt.errDetail)
			}
		})
	}
}

func TestWriteError_WithoutLogger(t *testing.T) {
	// Setup
	w := httptest.NewRecorder()

	// Execute
	WriteError(w, http.StatusBadRequest, "Test error", "test detail", nil)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "error")
	assert.Contains(t, w.Body.String(), "Test error")
	assert.Contains(t, w.Body.String(), "test detail")
}

func TestSuccessResponse_Structure(t *testing.T) {
	resp := SuccessResponse{
		Status:  "success",
		Message: "test message",
		Data:    "test data",
	}

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "test message", resp.Message)
	assert.Equal(t, "test data", resp.Data)
}

func TestErrorResponse_Structure(t *testing.T) {
	resp := ErrorResponse{
		Status:  "error",
		Message: "test error",
		Error:   "test detail",
	}

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "test error", resp.Message)
	assert.Equal(t, "test detail", resp.Error)
}
