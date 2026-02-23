package handler

import (
	"encoding/json"
	"net/http"

	"quickbite/config"
	"quickbite/internal/middleware"
	"quickbite/internal/model"
	"quickbite/internal/service"
	"quickbite/internal/utils"
)

type RestaurantHandler struct {
	cfg *config.Config
}

func NewRestaurantHandler(cfg *config.Config) *RestaurantHandler {
	return &RestaurantHandler{cfg: cfg}
}

func (h *RestaurantHandler) CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.CreateRestaurantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	restaurant, err := service.CreateRestaurant(&req, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, restaurant)
}

func (h *RestaurantHandler) GetRestaurantByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, "restaurant id is required")
		return
	}

	restaurant, err := service.GetRestaurantByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, restaurant)
}

func (h *RestaurantHandler) GetMyRestaurants(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	restaurants, err := service.GetRestaurantsByOwner(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to fetch restaurants")
		return
	}

	utils.WriteJSON(w, http.StatusOK, restaurants)
}

func (h *RestaurantHandler) GetAllRestaurants(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")

	restaurants, err := service.GetAllRestaurants(city)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to fetch restaurants")
		return
	}

	utils.WriteJSON(w, http.StatusOK, restaurants)
}

func (h *RestaurantHandler) UpdateRestaurant(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, "restaurant id is required")
		return
	}

	var req model.UpdateRestaurantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := service.UpdateRestaurant(id, &req, userID); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "restaurant updated successfully"})
}

func (h *RestaurantHandler) DeleteRestaurant(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, "restaurant id is required")
		return
	}

	if err := service.DeleteRestaurant(id, userID); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "restaurant deleted successfully"})
}
