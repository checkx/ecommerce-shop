package testutils

import (
	"fmt"
	"os"
	"testing"

	"ecommerce-shop/internal/config"
)

// TestConfigWithEnv creates a test configuration with environment variables
func TestConfigWithEnv() config.Config {
	return config.Config{
		Env:                   getEnvOrDefault("TEST_APP_ENV", "test"),
		HTTPAddr:              getEnvOrDefault("TEST_HTTP_ADDR", ":8080"),
		DBURL:                 getEnvOrDefault("TEST_DATABASE_URL", "postgres://test:test@localhost:5432/test_db?sslmode=disable"),
		JWTSecret:             getEnvOrDefault("TEST_JWT_SECRET", "test-secret-key-for-testing-only"),
		ReservationTTLMinutes: 15,
	}
}

// TestConfigWithDB creates a test configuration with specific database settings
func TestConfigWithDB(t *testing.T, host, port, user, password, name string) config.Config {
	return config.Config{
		Env:                   "test",
		HTTPAddr:              ":8080",
		DBURL:                 fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, name),
		JWTSecret:             "test-secret-key-for-testing-only",
		ReservationTTLMinutes: 15,
	}
}

// getEnvOrDefault returns the environment variable value or a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// SetTestEnv sets test environment variables
func SetTestEnv(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("JWT_SECRET", "test-secret-key")
	t.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/test_db?sslmode=disable")
	t.Setenv("HTTP_ADDR", ":8080")
}
