package testutils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"ecommerce-shop/internal/config"
)

// TestContext creates a test context
func TestContext() context.Context {
	return context.Background()
}

// TestGinContext creates a gin context for testing
func TestGinContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	return c, w
}

// TestGinContextWithBody creates a gin context with JSON body
func TestGinContextWithBody(t *testing.T, body interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	c, w := TestGinContext()

	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Failed to marshal body: %v", err)
	}

	c.Request.Body = http.NoBody
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = http.NoBody

	// Create a new request with the body
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	return c, w
}

// TestGinContextWithHeader creates a gin context with specific header
func TestGinContextWithHeader(t *testing.T, headerKey, headerValue string, body interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	c, w := TestGinContextWithBody(t, body)
	c.Request.Header.Set(headerKey, headerValue)
	return c, w
}

// TestValidator creates a validator for testing
func TestValidator() *validator.Validate {
	return validator.New()
}

// TestConfig creates a test configuration
func TestConfig() config.Config {
	return config.Config{
		Env:                   "test",
		HTTPAddr:              ":8080",
		DBURL:                 "postgres://test:test@localhost:5432/test_db?sslmode=disable",
		JWTSecret:             "test-secret-key",
		ReservationTTLMinutes: 15,
	}
}

// AssertJSONResponse checks if the response matches expected JSON
func AssertJSONResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedBody interface{}) {
	if w.Code != expectedStatus {
		t.Errorf("Expected status %d, got %d", expectedStatus, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
		return
	}

	expectedJSON, err := json.Marshal(expectedBody)
	if err != nil {
		t.Errorf("Failed to marshal expected body: %v", err)
		return
	}

	var expected map[string]interface{}
	if err := json.Unmarshal(expectedJSON, &expected); err != nil {
		t.Errorf("Failed to unmarshal expected body: %v", err)
		return
	}

	// Compare the response structure
	for key, expectedValue := range expected {
		if actualValue, exists := response[key]; !exists {
			t.Errorf("Expected key '%s' not found in response", key)
		} else if fmt.Sprintf("%v", actualValue) != fmt.Sprintf("%v", expectedValue) {
			t.Errorf("Expected %s to be %v, got %v", key, expectedValue, actualValue)
		}
	}
}

// AssertErrorResponse checks if the response is an error response
func AssertErrorResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedMessage string) {
	if w.Code != expectedStatus {
		t.Errorf("Expected status %d, got %d", expectedStatus, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
		return
	}

	if message, exists := response["message"]; !exists {
		t.Error("Expected 'message' field in error response")
	} else if message != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, message)
	}
}
