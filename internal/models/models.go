package models

import "time"

type User struct {
	ID           string    `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type Shop struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Warehouse struct {
	ID        string    `db:"id" json:"id"`
	ShopID    string    `db:"shop_id" json:"shop_id"`
	Name      string    `db:"name" json:"name"`
	Active    bool      `db:"active" json:"active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Product struct {
	ID        string    `db:"id" json:"id"`
	SKU       string    `db:"sku" json:"sku"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Inventory struct {
	WarehouseID string `db:"warehouse_id" json:"warehouse_id"`
	ProductID   string `db:"product_id" json:"product_id"`
	Quantity    int    `db:"quantity" json:"quantity"`
}

type Order struct {
	ID         string    `db:"id" json:"id"`
	UserID     string    `db:"user_id" json:"user_id"`
	ShopID     string    `db:"shop_id" json:"shop_id"`
	Status     string    `db:"status" json:"status"`
	TotalCents int64     `db:"total_cents" json:"total_cents"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

type OrderItem struct {
	OrderID   string `db:"order_id" json:"order_id"`
	ProductID string `db:"product_id" json:"product_id"`
	Quantity  int    `db:"quantity" json:"quantity"`
}

type Reservation struct {
	ID          string    `db:"id" json:"id"`
	OrderID     string    `db:"order_id" json:"order_id"`
	WarehouseID string    `db:"warehouse_id" json:"warehouse_id"`
	ProductID   string    `db:"product_id" json:"product_id"`
	Quantity    int       `db:"quantity" json:"quantity"`
	ExpiresAt   time.Time `db:"expires_at" json:"expires_at"`
	Released    bool      `db:"released" json:"released"`
}
