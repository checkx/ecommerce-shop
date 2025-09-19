package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransferReq_Structure(t *testing.T) {
	req := TransferReq{
		From:      "wh-1",
		To:        "wh-2",
		ProductID: "prod-1",
		Quantity:  5,
	}

	assert.Equal(t, "wh-1", req.From)
	assert.Equal(t, "wh-2", req.To)
	assert.Equal(t, "prod-1", req.ProductID)
	assert.Equal(t, 5, req.Quantity)
}

func TestTransferResponse_Structure(t *testing.T) {
	response := TransferResponse{
		From:      "wh-1",
		To:        "wh-2",
		ProductID: "prod-1",
		Quantity:  5,
		Status:    "ok",
	}

	assert.Equal(t, "wh-1", response.From)
	assert.Equal(t, "wh-2", response.To)
	assert.Equal(t, "prod-1", response.ProductID)
	assert.Equal(t, 5, response.Quantity)
	assert.Equal(t, "ok", response.Status)
}
