package service

import (
	"errors"
	"quickbite/internal/model"
	"quickbite/internal/repository"
)

// ====== MENU CATEGORIES ======

func CreateCategory(req *model.CreateCategoryRequest, userID string) (*model.MenuCategory, error) {
	if req.Name == "" {
		return nil, errors.New("category name is required")
	}

	restaurant, err := repository.GetRestaurantByID(req.RestaurantID)
	if err != nil {
		return nil, errors.New("restaurant not found")
	}

	if restaurant.OwnerID != userID {
		return nil, errors.New("unauthorized: you don't own this restaurant")
	}

	category := &model.MenuCategory{
		RestaurantID: req.RestaurantID,
		Name:         req.Name,
		DisplayOrder: req.DisplayOrder,
	}

	if err := repository.CreateCategory(category); err != nil {
		return nil, errors.New("failed to create category")
	}

	return category, nil
}

func GetCategoriesByRestaurant(restaurantID string) ([]model.MenuCategory, error) {
	return repository.GetCategoriesByRestaurant(restaurantID)
}

func DeleteCategory(id string, userID string) error {
	// Get the category to find which restaurant it belongs to
	category, err := repository.GetCategoryByID(id)
	if err != nil {
		return errors.New("category not found")
	}

	// Get the restaurant to verify ownership
	restaurant, err := repository.GetRestaurantByID(category.RestaurantID)
	if err != nil {
		return errors.New("restaurant not found")
	}

	// Verify user owns the restaurant
	if restaurant.OwnerID != userID {
		return errors.New("unauthorized: you don't own this restaurant")
	}

	return repository.DeleteCategory(id)
}

// ====== MENU ITEMS ======

func CreateMenuItem(req *model.CreateMenuItemRequest, userID string) (*model.MenuItem, error) {
	if req.Name == "" || req.Price <= 0 {
		return nil, errors.New("name and valid price are required")
	}

	// Get the category to find which restaurant it belongs to
	category, err := repository.GetCategoryByID(req.CategoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	// Get the restaurant to verify ownership
	restaurant, err := repository.GetRestaurantByID(category.RestaurantID)
	if err != nil {
		return nil, errors.New("restaurant not found")
	}

	// Verify user owns the restaurant
	if restaurant.OwnerID != userID {
		return nil, errors.New("unauthorized: you don't own this restaurant")
	}

	item := &model.MenuItem{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		ImageURL:    req.ImageURL,
		IsAvailable: true,
		IsVeg:       req.IsVeg,
	}

	if err := repository.CreateMenuItem(item); err != nil {
		return nil, errors.New("failed to create menu item")
	}

	return item, nil
}

func GetMenuItemsByCategory(categoryID string) ([]model.MenuItem, error) {
	return repository.GetMenuItemsByCategory(categoryID)
}

func UpdateMenuItem(id string, req *model.UpdateMenuItemRequest, userID string) error {
	if req.Name == "" || req.Price <= 0 {
		return errors.New("name and valid price are required")
	}

	// Get the item
	item, err := repository.GetMenuItemByID(id)
	if err != nil {
		return errors.New("menu item not found")
	}

	// Get the category
	category, err := repository.GetCategoryByID(item.CategoryID)
	if err != nil {
		return errors.New("category not found")
	}

	// Get the restaurant
	restaurant, err := repository.GetRestaurantByID(category.RestaurantID)
	if err != nil {
		return errors.New("restaurant not found")
	}

	// Verify ownership
	if restaurant.OwnerID != userID {
		return errors.New("unauthorized: you don't own this restaurant")
	}

	return repository.UpdateMenuItem(id, req)
}

func DeleteMenuItem(id string, userID string) error {
	// Get the item
	item, err := repository.GetMenuItemByID(id)
	if err != nil {
		return errors.New("menu item not found")
	}

	// Get the category
	category, err := repository.GetCategoryByID(item.CategoryID)
	if err != nil {
		return errors.New("category not found")
	}

	// Get the restaurant
	restaurant, err := repository.GetRestaurantByID(category.RestaurantID)
	if err != nil {
		return errors.New("restaurant not found")
	}

	// Verify ownership
	if restaurant.OwnerID != userID {
		return errors.New("unauthorized: you don't own this restaurant")
	}

	return repository.DeleteMenuItem(id)
}
