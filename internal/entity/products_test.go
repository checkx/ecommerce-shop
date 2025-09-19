package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductResponse_Structure(t *testing.T) {
	response := ProductResponse{
		ID:        "prod-123",
		SKU:       "SKU001",
		Name:      "Test Product",
		Available: 10,
	}

	assert.Equal(t, "prod-123", response.ID)
	assert.Equal(t, "SKU001", response.SKU)
	assert.Equal(t, "Test Product", response.Name)
	assert.Equal(t, 10, response.Available)
}
