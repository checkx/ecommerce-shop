package entity

type TransferReq struct {
	From      string `json:"from" binding:"required"`
	To        string `json:"to" binding:"required"`
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

type TransferResponse struct {
	From      string `json:"from"`
	To        string `json:"to"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Status    string `json:"status"`
}
