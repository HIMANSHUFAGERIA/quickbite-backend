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

type MenuHandler struct {
	cfg *config.Config
}

func NewMenuHandler(cfg *config.Config) *MenuHandler {
	return &MenuHandler{cfg: cfg}
}

// ====== MENU CATEGORIES ======

func (h *MenuHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	category, err := service.CreateCategory(&req, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, category)
}

func (h *MenuHandler) GetCategoriesByRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurantID := r.PathValue("restaurant_id")
	if restaurantID == "" {
		utils.WriteError(w, http.StatusBadRequest, "restaurant id is required")
		return
	}

	categories, err := service.GetCategoriesByRestaurant(restaurantID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to fetch categories")
		return
	}

	utils.WriteJSON(w, http.StatusOK, categories)
}

func (h *MenuHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		log.Println("DeleteCategory: unauthorized - no userID in context")
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		log.Println("DeleteCategory: category id is required")
		utils.WriteError(w, http.StatusBadRequest, "category id is required")
		return
	}

	log.Printf("DeleteCategory: attempting to delete category %s by user %s", id, userID)

	if err := service.DeleteCategory(id, userID); err != nil {
		log.Printf("DeleteCategory ERROR: %v", err)
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("DeleteCategory: successfully deleted category %s", id)
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "category deleted successfully"})
}

// ====== MENU ITEMS ======

func (h *MenuHandler) CreateMenuItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.CreateMenuItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := service.CreateMenuItem(&req, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, item)
}

func (h *MenuHandler) GetMenuItemsByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := r.PathValue("category_id")
	if categoryID == "" {
		utils.WriteError(w, http.StatusBadRequest, "category id is required")
		return
	}

	items, err := service.GetMenuItemsByCategory(categoryID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to fetch menu items")
		return
	}

	utils.WriteJSON(w, http.StatusOK, items)
}

func (h *MenuHandler) UpdateMenuItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, "menu item id is required")
		return
	}

	var req model.UpdateMenuItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := service.UpdateMenuItem(id, &req, userID); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "menu item updated successfully"})
}

func (h *MenuHandler) DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, "menu item id is required")
		return
	}

	if err := service.DeleteMenuItem(id, userID); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "menu item deleted successfully"})
}
