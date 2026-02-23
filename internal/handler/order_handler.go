package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"quickbite/config"
	"quickbite/internal/middleware"
	"quickbite/internal/model"
	"quickbite/internal/service"
	"quickbite/internal/utils"
)

type OrderHandler struct {
	cfg *config.Config
}

func NewOrderHandler(cfg *config.Config) *OrderHandler {
	return &OrderHandler{cfg: cfg}
}

// CreateOrder handles POST /api/orders
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	log.Printf("CreateOrder: user %s ordering from restaurant %s", userID, req.RestaurantID)

	order, err := service.CreateOrder(&req, userID)
	if err != nil {
		log.Printf("CreateOrder error: %v", err)
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("CreateOrder: order %s created successfully", order.ID)
	utils.WriteJSON(w, http.StatusCreated, order)
}

// GetOrderByID handles GET /api/orders/:id
func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	orderID := r.PathValue("id")
	if orderID == "" {
		utils.WriteError(w, http.StatusBadRequest, "order id is required")
		return
	}

	order, err := service.GetOrderByID(orderID, userID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, order)
}

// GetMyOrders handles GET /api/orders/my/list
func (h *OrderHandler) GetMyOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	orders, err := service.GetMyOrders(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to fetch orders")
		return
	}

	utils.WriteJSON(w, http.StatusOK, orders)
}

// GetRestaurantOrders handles GET /api/restaurants/:id/orders
func (h *OrderHandler) GetRestaurantOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	restaurantID := r.PathValue("id")
	if restaurantID == "" {
		utils.WriteError(w, http.StatusBadRequest, "restaurant id is required")
		return
	}

	orders, err := service.GetRestaurantOrders(restaurantID, userID)
	if err != nil {
		log.Printf("GetRestaurantOrders error: %v", err)
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, orders)
}

// UpdateOrderStatus handles PUT /api/orders/:id/status
func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	orderID := r.PathValue("id")
	if orderID == "" {
		utils.WriteError(w, http.StatusBadRequest, "order id is required")
		return
	}

	var body struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.Status == "" {
		utils.WriteError(w, http.StatusBadRequest, "status is required")
		return
	}

	log.Printf("UpdateOrderStatus: order %s to status %s by user %s", orderID, body.Status, userID)

	if err := service.UpdateOrderStatus(orderID, body.Status, userID); err != nil {
		log.Printf("UpdateOrderStatus error: %v", err)
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "order status updated successfully"})
}

// CancelOrder handles POST /api/orders/:id/cancel
func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	orderID := r.PathValue("id")
	if orderID == "" {
		utils.WriteError(w, http.StatusBadRequest, "order id is required")
		return
	}

	log.Printf("CancelOrder: user %s cancelling order %s", userID, orderID)

	if err := service.CancelOrder(orderID, userID); err != nil {
		log.Printf("CancelOrder error: %v", err)
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "order cancelled successfully"})
}
