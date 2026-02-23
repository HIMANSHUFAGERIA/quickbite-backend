package model

import "time"

type Order struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	RestaurantID    string    `json:"restaurant_id"`
	Status          string    `json:"status"`
	TotalAmount     float64   `json:"total_amount"`
	DeliveryFee     float64   `json:"delivery_fee"`
	DeliveryAddress string    `json:"delivery_address"`
	PaymentMethod   string    `json:"payment_method"`
	PaymentStatus   string    `json:"payment_status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type OrderItem struct {
	ID         string    `json:"id"`
	OrderID    string    `json:"order_id"`
	MenuItemID string    `json:"menu_item_id"`
	Quantity   int       `json:"quantity"`
	Price      float64   `json:"price"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateOrderRequest struct {
	RestaurantID    string           `json:"restaurant_id"`
	Items           []OrderItemInput `json:"items"`
	DeliveryAddress string           `json:"delivery_address"`
	PaymentMethod   string           `json:"payment_method"`
}

type OrderItemInput struct {
	MenuItemID string `json:"menu_item_id"`
	Quantity   int    `json:"quantity"`
}

type OrderWithDetails struct {
	Order
	RestaurantName string              `json:"restaurant_name"`
	Items          []OrderItemWithMenu `json:"items"`
}

type OrderItemWithMenu struct {
	OrderItem
	ItemName  string `json:"item_name"`
	ItemImage string `json:"item_image"`
	IsVeg     bool   `json:"is_veg"`
}
