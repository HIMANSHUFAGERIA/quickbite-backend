package repository

import (
	"context"
	"quickbite/db"
	"quickbite/internal/model"
)

// CreateOrder inserts a new order and returns it with generated ID
func CreateOrder(order *model.Order) error {
	query := `
		INSERT INTO orders (user_id, restaurant_id, status, total_amount, delivery_fee, delivery_address, payment_method, payment_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	return db.DB.QueryRow(
		context.Background(),
		query,
		order.UserID,
		order.RestaurantID,
		order.Status,
		order.TotalAmount,
		order.DeliveryFee,
		order.DeliveryAddress,
		order.PaymentMethod,
		order.PaymentStatus,
	).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
}

// CreateOrderItem inserts an order item
func CreateOrderItem(item *model.OrderItem) error {
	query := `
		INSERT INTO order_items (order_id, menu_item_id, quantity, price)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	return db.DB.QueryRow(
		context.Background(),
		query,
		item.OrderID,
		item.MenuItemID,
		item.Quantity,
		item.Price,
	).Scan(&item.ID, &item.CreatedAt)
}

// GetOrderByID fetches a single order by ID
func GetOrderByID(id string) (*model.Order, error) {
	query := `
		SELECT id, user_id, restaurant_id, status, total_amount, delivery_fee, 
		       delivery_address, payment_method, payment_status, created_at, updated_at
		FROM orders
		WHERE id = $1::uuid
	`

	order := &model.Order{}

	err := db.DB.QueryRow(context.Background(), query, id).Scan(
		&order.ID,
		&order.UserID,
		&order.RestaurantID,
		&order.Status,
		&order.TotalAmount,
		&order.DeliveryFee,
		&order.DeliveryAddress,
		&order.PaymentMethod,
		&order.PaymentStatus,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// GetOrderWithDetails fetches order with restaurant name and items with menu details
func GetOrderWithDetails(id string) (*model.OrderWithDetails, error) {
	// First get the order
	order, err := GetOrderByID(id)
	if err != nil {
		return nil, err
	}

	// Get restaurant name
	restaurantQuery := `SELECT name FROM restaurants WHERE id = $1::uuid`
	var restaurantName string
	err = db.DB.QueryRow(context.Background(), restaurantQuery, order.RestaurantID).Scan(&restaurantName)
	if err != nil {
		return nil, err
	}

	// Get order items with menu details
	itemsQuery := `
		SELECT 
			oi.id, oi.order_id, oi.menu_item_id, oi.quantity, oi.price, oi.created_at,
			mi.name, mi.image_url, mi.is_veg
		FROM order_items oi
		JOIN menu_items mi ON oi.menu_item_id = mi.id
		WHERE oi.order_id = $1::uuid
	`

	rows, err := db.DB.Query(context.Background(), itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.OrderItemWithMenu

	for rows.Next() {
		var item model.OrderItemWithMenu
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.MenuItemID,
			&item.Quantity,
			&item.Price,
			&item.CreatedAt,
			&item.ItemName,
			&item.ItemImage,
			&item.IsVeg,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return &model.OrderWithDetails{
		Order:          *order,
		RestaurantName: restaurantName,
		Items:          items,
	}, nil
}

// GetOrdersByUser fetches all orders for a specific user
func GetOrdersByUser(userID string) ([]model.OrderWithDetails, error) {
	query := `
		SELECT o.id, o.user_id, o.restaurant_id, o.status, o.total_amount, 
		       o.delivery_fee, o.delivery_address, o.payment_method, o.payment_status,
		       o.created_at, o.updated_at, r.name as restaurant_name
		FROM orders o
		JOIN restaurants r ON o.restaurant_id = r.id
		WHERE o.user_id = $1::uuid
		ORDER BY o.created_at DESC
	`

	rows, err := db.DB.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.OrderWithDetails

	for rows.Next() {
		var orderDetail model.OrderWithDetails
		err := rows.Scan(
			&orderDetail.ID,
			&orderDetail.UserID,
			&orderDetail.RestaurantID,
			&orderDetail.Status,
			&orderDetail.TotalAmount,
			&orderDetail.DeliveryFee,
			&orderDetail.DeliveryAddress,
			&orderDetail.PaymentMethod,
			&orderDetail.PaymentStatus,
			&orderDetail.CreatedAt,
			&orderDetail.UpdatedAt,
			&orderDetail.RestaurantName,
		)
		if err != nil {
			return nil, err
		}

		// Get items for this order
		items, err := getOrderItems(orderDetail.ID)
		if err != nil {
			return nil, err
		}
		orderDetail.Items = items

		orders = append(orders, orderDetail)
	}

	return orders, nil
}

// GetOrdersByRestaurant fetches all orders for a specific restaurant
func GetOrdersByRestaurant(restaurantID string) ([]model.OrderWithDetails, error) {
	query := `
		SELECT o.id, o.user_id, o.restaurant_id, o.status, o.total_amount, 
		       o.delivery_fee, o.delivery_address, o.payment_method, o.payment_status,
		       o.created_at, o.updated_at, r.name as restaurant_name
		FROM orders o
		JOIN restaurants r ON o.restaurant_id = r.id
		WHERE o.restaurant_id = $1::uuid
		ORDER BY o.created_at DESC
	`

	rows, err := db.DB.Query(context.Background(), query, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.OrderWithDetails

	for rows.Next() {
		var orderDetail model.OrderWithDetails
		err := rows.Scan(
			&orderDetail.ID,
			&orderDetail.UserID,
			&orderDetail.RestaurantID,
			&orderDetail.Status,
			&orderDetail.TotalAmount,
			&orderDetail.DeliveryFee,
			&orderDetail.DeliveryAddress,
			&orderDetail.PaymentMethod,
			&orderDetail.PaymentStatus,
			&orderDetail.CreatedAt,
			&orderDetail.UpdatedAt,
			&orderDetail.RestaurantName,
		)
		if err != nil {
			return nil, err
		}

		// Get items for this order
		items, err := getOrderItems(orderDetail.ID)
		if err != nil {
			return nil, err
		}
		orderDetail.Items = items

		orders = append(orders, orderDetail)
	}

	return orders, nil
}

// UpdateOrderStatus updates only the status of an order
func UpdateOrderStatus(orderID string, status string) error {
	query := `
		UPDATE orders
		SET status = $1, updated_at = NOW()
		WHERE id = $2::uuid
	`

	_, err := db.DB.Exec(context.Background(), query, status, orderID)
	return err
}

// Helper function to get order items with menu details
func getOrderItems(orderID string) ([]model.OrderItemWithMenu, error) {
	query := `
		SELECT 
			oi.id, oi.order_id, oi.menu_item_id, oi.quantity, oi.price, oi.created_at,
			mi.name, mi.image_url, mi.is_veg
		FROM order_items oi
		JOIN menu_items mi ON oi.menu_item_id = mi.id
		WHERE oi.order_id = $1::uuid
	`

	rows, err := db.DB.Query(context.Background(), query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.OrderItemWithMenu

	for rows.Next() {
		var item model.OrderItemWithMenu
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.MenuItemID,
			&item.Quantity,
			&item.Price,
			&item.CreatedAt,
			&item.ItemName,
			&item.ItemImage,
			&item.IsVeg,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
