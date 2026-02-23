package service

import (
	"errors"
	"quickbite/internal/model"
	"quickbite/internal/repository"
)

// CreateOrder validates items, calculates total, and creates order with items
func CreateOrder(req *model.CreateOrderRequest, userID string) (*model.OrderWithDetails, error) {
	// Validate request
	if req.RestaurantID == "" {
		return nil, errors.New("restaurant_id is required")
	}
	if len(req.Items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}
	if req.DeliveryAddress == "" {
		return nil, errors.New("delivery_address is required")
	}
	if req.PaymentMethod == "" {
		return nil, errors.New("payment_method is required")
	}

	// Validate restaurant exists and is active
	restaurant, err := repository.GetRestaurantByID(req.RestaurantID)
	if err != nil {
		return nil, errors.New("restaurant not found")
	}
	if !restaurant.IsActive {
		return nil, errors.New("restaurant is currently closed")
	}

	// Validate items and calculate total
	var totalAmount float64
	var validatedItems []model.OrderItem

	for _, itemInput := range req.Items {
		if itemInput.Quantity <= 0 {
			return nil, errors.New("item quantity must be greater than 0")
		}

		// Get menu item to verify it exists and is available
		menuItem, err := repository.GetMenuItemByID(itemInput.MenuItemID)
		if err != nil {
			return nil, errors.New("menu item not found: " + itemInput.MenuItemID)
		}

		if !menuItem.IsAvailable {
			return nil, errors.New("item is not available: " + menuItem.Name)
		}

		// Calculate item total
		itemTotal := menuItem.Price * float64(itemInput.Quantity)
		totalAmount += itemTotal

		// Store validated item
		validatedItems = append(validatedItems, model.OrderItem{
			MenuItemID: itemInput.MenuItemID,
			Quantity:   itemInput.Quantity,
			Price:      menuItem.Price, // Store current price
		})
	}

	// Set delivery fee (fixed for now, can be dynamic later)
	deliveryFee := 50.0

	// Create order
	order := &model.Order{
		UserID:          userID,
		RestaurantID:    req.RestaurantID,
		Status:          "pending",
		TotalAmount:     totalAmount,
		DeliveryFee:     deliveryFee,
		DeliveryAddress: req.DeliveryAddress,
		PaymentMethod:   req.PaymentMethod,
		PaymentStatus:   "pending", // Will change to "paid" after payment gateway integration
	}

	// Insert order into database
	if err := repository.CreateOrder(order); err != nil {
		return nil, errors.New("failed to create order")
	}

	// Insert order items
	for i := range validatedItems {
		validatedItems[i].OrderID = order.ID
		if err := repository.CreateOrderItem(&validatedItems[i]); err != nil {
			// TODO: In production, wrap this in a transaction and rollback if items fail
			return nil, errors.New("failed to create order items")
		}
	}

	// Fetch and return complete order details
	orderDetails, err := repository.GetOrderWithDetails(order.ID)
	if err != nil {
		return nil, errors.New("order created but failed to fetch details")
	}

	return orderDetails, nil
}

// GetOrderByID fetches order details for a user
func GetOrderByID(orderID string, userID string) (*model.OrderWithDetails, error) {
	orderDetails, err := repository.GetOrderWithDetails(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	// Verify user owns this order
	if orderDetails.UserID != userID {
		return nil, errors.New("unauthorized: you don't own this order")
	}

	return orderDetails, nil
}

// GetMyOrders fetches all orders for a user
func GetMyOrders(userID string) ([]model.OrderWithDetails, error) {
	return repository.GetOrdersByUser(userID)
}

// GetRestaurantOrders fetches all orders for a restaurant (owner only)
func GetRestaurantOrders(restaurantID string, userID string) ([]model.OrderWithDetails, error) {
	// Verify user owns the restaurant
	restaurant, err := repository.GetRestaurantByID(restaurantID)
	if err != nil {
		return nil, errors.New("restaurant not found")
	}

	if restaurant.OwnerID != userID {
		return nil, errors.New("unauthorized: you don't own this restaurant")
	}

	return repository.GetOrdersByRestaurant(restaurantID)
}

// UpdateOrderStatus allows restaurant owners to update order status
func UpdateOrderStatus(orderID string, newStatus string, userID string) error {
	// Validate status
	validStatuses := map[string]bool{
		"pending":          true,
		"confirmed":        true,
		"preparing":        true,
		"ready":            true,
		"out_for_delivery": true,
		"delivered":        true,
		"cancelled":        true,
	}

	if !validStatuses[newStatus] {
		return errors.New("invalid order status")
	}

	// Get order to verify ownership
	order, err := repository.GetOrderByID(orderID)
	if err != nil {
		return errors.New("order not found")
	}

	// Verify user owns the restaurant this order belongs to
	restaurant, err := repository.GetRestaurantByID(order.RestaurantID)
	if err != nil {
		return errors.New("restaurant not found")
	}

	if restaurant.OwnerID != userID {
		return errors.New("unauthorized: you don't own this restaurant")
	}

	// Validate status transitions (state machine logic)
	// For now, allow any transition, but in production you'd enforce rules like:
	// - Can't go from "delivered" to "pending"
	// - Can't change status after "delivered" or "cancelled"
	if order.Status == "delivered" || order.Status == "cancelled" {
		return errors.New("cannot update status of completed order")
	}

	return repository.UpdateOrderStatus(orderID, newStatus)
}

// CancelOrder allows customers to cancel their order (only if status is pending or confirmed)
func CancelOrder(orderID string, userID string) error {
	order, err := repository.GetOrderByID(orderID)
	if err != nil {
		return errors.New("order not found")
	}

	// Verify user owns this order
	if order.UserID != userID {
		return errors.New("unauthorized: you don't own this order")
	}

	// Only allow cancellation if order is pending or confirmed
	if order.Status != "pending" && order.Status != "confirmed" {
		return errors.New("cannot cancel order in current status")
	}

	return repository.UpdateOrderStatus(orderID, "cancelled")
}
