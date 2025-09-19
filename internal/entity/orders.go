package entity

type OrderItemReq struct {
	ProductID string `json:"product_id" validate:"required,uuid"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}
type CreateOrderReq struct {
	ShopID string         `json:"shop_id" validate:"required,uuid"`
	Items  []OrderItemReq `json:"items" validate:"required,min=1,dive"`
}

type OrderResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}
