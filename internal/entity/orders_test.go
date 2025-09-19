package entity

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestOrderItemReq_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name    string
		req     OrderItemReq
		wantErr bool
	}{
		{
			name: "valid order item",
			req: OrderItemReq{
				ProductID: "550e8400-e29b-41d4-a716-446655440000",
				Quantity:  5,
			},
			wantErr: false,
		},
		{
			name: "missing product_id",
			req: OrderItemReq{
				Quantity: 5,
			},
			wantErr: true,
		},
		{
			name: "invalid product_id format",
			req: OrderItemReq{
				ProductID: "invalid-uuid",
				Quantity:  5,
			},
			wantErr: true,
		},
		{
			name: "missing quantity",
			req: OrderItemReq{
				ProductID: "550e8400-e29b-41d4-a716-446655440000",
			},
			wantErr: true,
		},
		{
			name: "quantity zero",
			req: OrderItemReq{
				ProductID: "550e8400-e29b-41d4-a716-446655440000",
				Quantity:  0,
			},
			wantErr: true,
		},
		{
			name: "quantity negative",
			req: OrderItemReq{
				ProductID: "550e8400-e29b-41d4-a716-446655440000",
				Quantity:  -1,
			},
			wantErr: true,
		},
		{
			name: "quantity one",
			req: OrderItemReq{
				ProductID: "550e8400-e29b-41d4-a716-446655440000",
				Quantity:  1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateOrderReq_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name    string
		req     CreateOrderReq
		wantErr bool
	}{
		{
			name: "valid order request",
			req: CreateOrderReq{
				ShopID: "550e8400-e29b-41d4-a716-446655440000",
				Items: []OrderItemReq{
					{
						ProductID: "550e8400-e29b-41d4-a716-446655440001",
						Quantity:  2,
					},
					{
						ProductID: "550e8400-e29b-41d4-a716-446655440002",
						Quantity:  1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing shop_id",
			req: CreateOrderReq{
				Items: []OrderItemReq{
					{
						ProductID: "550e8400-e29b-41d4-a716-446655440001",
						Quantity:  2,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid shop_id format",
			req: CreateOrderReq{
				ShopID: "invalid-uuid",
				Items: []OrderItemReq{
					{
						ProductID: "550e8400-e29b-41d4-a716-446655440001",
						Quantity:  2,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing items",
			req: CreateOrderReq{
				ShopID: "550e8400-e29b-41d4-a716-446655440000",
			},
			wantErr: true,
		},
		{
			name: "empty items",
			req: CreateOrderReq{
				ShopID: "550e8400-e29b-41d4-a716-446655440000",
				Items:  []OrderItemReq{},
			},
			wantErr: true,
		},
		{
			name: "single item",
			req: CreateOrderReq{
				ShopID: "550e8400-e29b-41d4-a716-446655440000",
				Items: []OrderItemReq{
					{
						ProductID: "550e8400-e29b-41d4-a716-446655440001",
						Quantity:  1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid item in items",
			req: CreateOrderReq{
				ShopID: "550e8400-e29b-41d4-a716-446655440000",
				Items: []OrderItemReq{
					{
						ProductID: "invalid-uuid",
						Quantity:  1,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOrderResponse_Structure(t *testing.T) {
	response := OrderResponse{
		ID:     "order-123",
		Status: "reserved",
	}

	assert.Equal(t, "order-123", response.ID)
	assert.Equal(t, "reserved", response.Status)
}
