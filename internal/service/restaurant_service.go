package service

import (
	"errors"
	"quickbite/internal/model"
	"quickbite/internal/repository"
)

func CreateRestaurant(req *model.CreateRestaurantRequest, ownerID string) (*model.Restaurant, error) {
	if req.Name == "" || req.Address == "" || req.City == "" {
		return nil, errors.New("name, address and city are required")
	}

	restaurant := &model.Restaurant{
		OwnerID:     ownerID,
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		City:        req.City,
		ImageURL:    req.ImageURL,
		IsActive:    true,
		Rating:      0.0,
	}

	if err := repository.CreateRestaurant(restaurant); err != nil {
		return nil, errors.New("failed to create restaurant")
	}

	return restaurant, nil
}

func GetRestaurantByID(id string) (*model.Restaurant, error) {
	restaurant, err := repository.GetRestaurantByID(id)
	if err != nil {
		return nil, errors.New("restaurant not found")
	}
	return restaurant, nil
}

func GetRestaurantsByOwner(ownerID string) ([]model.Restaurant, error) {
	return repository.GetRestaurantsByOwner(ownerID)
}

func GetAllRestaurants(city string) ([]model.Restaurant, error) {
	return repository.GetAllRestaurants(city)
}

func UpdateRestaurant(id string, req *model.UpdateRestaurantRequest, userID string) error {
	restaurant, err := repository.GetRestaurantByID(id)
	if err != nil {
		return errors.New("restaurant not found")
	}

	if restaurant.OwnerID != userID {
		return errors.New("unauthorized: you don't own this restaurant")
	}

	if req.Name == "" || req.Address == "" || req.City == "" {
		return errors.New("name, address and city are required")
	}

	return repository.UpdateRestaurant(id, req)
}

func DeleteRestaurant(id string, userID string) error {
	restaurant, err := repository.GetRestaurantByID(id)
	if err != nil {
		return errors.New("restaurant not found")
	}

	if restaurant.OwnerID != userID {
		return errors.New("unauthorized: you don't own this restaurant")
	}

	return repository.DeleteRestaurant(id)
}
