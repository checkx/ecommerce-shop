package entity

type ProductResponse struct {
	ID        string `json:"id"`
	SKU       string `json:"sku"`
	Name      string `json:"name"`
	Available int    `json:"available"`
}
