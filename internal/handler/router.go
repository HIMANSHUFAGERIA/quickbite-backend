package handler

import (
	"net/http"

	"quickbite/config"
	"quickbite/internal/middleware"
	"quickbite/internal/utils"
)

func NewRouter(cfg *config.Config) *http.ServeMux {
	mux := http.NewServeMux()

	// Initialize handlers
	authHandler := NewAuthHandler(cfg)
	restaurantHandler := NewRestaurantHandler(cfg)
	menuHandler := NewMenuHandler(cfg)
	orderHandler := NewOrderHandler(cfg)

	// Health check
	mux.HandleFunc("GET /health", healthCheck)

	// ====== AUTH ROUTES (Public) ======
	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)

	// ====== RESTAURANT ROUTES ======

	// Public routes - no auth required
	mux.HandleFunc("GET /api/restaurants", restaurantHandler.GetAllRestaurants)
	mux.HandleFunc("GET /api/restaurants/{id}", restaurantHandler.GetRestaurantByID)

	// Protected routes - require auth
	mux.Handle("GET /api/restaurants/my/list",
		middleware.Auth(cfg)(
			middleware.RequireRole("restaurant_owner")(
				http.HandlerFunc(restaurantHandler.GetMyRestaurants),
			),
		),
	)

	mux.Handle("POST /api/restaurants",
		middleware.Auth(cfg)(
			middleware.RequireRole("restaurant_owner")(
				http.HandlerFunc(restaurantHandler.CreateRestaurant),
			),
		),
	)

	mux.Handle("PUT /api/restaurants/{id}",
		middleware.Auth(cfg)(
			middleware.RequireRole("restaurant_owner")(
				http.HandlerFunc(restaurantHandler.UpdateRestaurant),
			),
		),
	)

	mux.Handle("DELETE /api/restaurants/{id}",
		middleware.Auth(cfg)(
			middleware.RequireRole("restaurant_owner")(
				http.HandlerFunc(restaurantHandler.DeleteRestaurant),
			),
		),
	)

	// ====== MENU CATEGORY ROUTES ======

	// Public route
	mux.HandleFunc("GET /api/restaurants/{restaurant_id}/categories", menuHandler.GetCategoriesByRestaurant)

	// Protected routes - require auth and restaurant_owner role
	mux.Handle("POST /api/menu/categories",
		middleware.Auth(cfg)(
			middleware.RequireRole("restaurant_owner")(
				http.HandlerFunc(menuHandler.CreateCategory),
			),
		),
	)

	mux.Handle("DELETE /api/menu/categories/{id}",
		middleware.Auth(cfg)(
			middleware.RequireRole("restaurant_owner")(
				http.HandlerFunc(menuHandler.DeleteCategory),
			),
		),
	)

	// ====== MENU ITEM ROUTES ======

	// Public route
	mux.HandleFunc("GET /api/categories/{category_id}/items", menuHandler.GetMenuItemsByCategory)

	// Protected routes - require auth and restaurant_owner role
	mux.Handle("POST /api/menu/items",
		middleware.Auth(cfg)(
			middleware.RequireRole("restaurant_owner")(
				http.HandlerFunc(menuHandler.CreateMenuItem),
			),
		),
	)

	mux.Handle("PUT /api/menu/items/{id}",
		middleware.Auth(cfg)(
			middleware.RequireRole("restaurant_owner")(
				http.HandlerFunc(menuHandler.UpdateMenuItem),
			),
		),
	)

	mux.Handle("DELETE /api/menu/items/{id}",
		middleware.Auth(cfg)(
			middleware.RequireRole("restaurant_owner")(
				http.HandlerFunc(menuHandler.DeleteMenuItem),
			),
		),
	)

	// ====== ORDER ROUTES ======

	// Customer routes - require auth (any authenticated user)
	mux.Handle("POST /api/orders",
		middleware.Auth(cfg)(
			http.HandlerFunc(orderHandler.CreateOrder),
		),
	)

	mux.Handle("GET /api/orders/my/list",
		middleware.Auth(cfg)(
			http.HandlerFunc(orderHandler.GetMyOrders),
		),
	)

	mux.Handle("GET /api/orders/{id}",
		middleware.Auth(cfg)(
			http.HandlerFunc(orderHandler.GetOrderByID),
		),
	)

	mux.Handle("POST /api/orders/{id}/cancel",
		middleware.Auth(cfg)(
			http.HandlerFunc(orderHandler.CancelOrder),
		),
	)

	// Restaurant owner routes - require auth and restaurant_owner role
	mux.Handle("GET /api/restaurants/{id}/orders",
		middleware.Auth(cfg)(
			middleware.RequireRole("restaurant_owner")(
				http.HandlerFunc(orderHandler.GetRestaurantOrders),
			),
		),
	)

	mux.Handle("PUT /api/orders/{id}/status",
		middleware.Auth(cfg)(
			middleware.RequireRole("restaurant_owner")(
				http.HandlerFunc(orderHandler.UpdateOrderStatus),
			),
		),
	)

	return mux
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"app":    "QuickBite",
	})
}
