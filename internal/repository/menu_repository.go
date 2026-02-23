package repository

import (
	"context"
	"quickbite/db"
	"quickbite/internal/model"
)

// ====== MENU CATEGORIES ======

func CreateCategory(category *model.MenuCategory) error {
	query := `
		INSERT INTO menu_categories (restaurant_id, name, display_order)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	return db.DB.QueryRow(
		context.Background(),
		query,
		category.RestaurantID,
		category.Name,
		category.DisplayOrder,
	).Scan(&category.ID, &category.CreatedAt)
}

func GetCategoriesByRestaurant(restaurantID string) ([]model.MenuCategory, error) {
	query := `
		SELECT id, restaurant_id, name, display_order, created_at
		FROM menu_categories
		WHERE restaurant_id = $1
		ORDER BY display_order ASC
	`

	rows, err := db.DB.Query(context.Background(), query, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.MenuCategory

	for rows.Next() {
		var c model.MenuCategory
		err := rows.Scan(
			&c.ID,
			&c.RestaurantID,
			&c.Name,
			&c.DisplayOrder,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func GetCategoryByID(id string) (*model.MenuCategory, error) {
	query := `
		SELECT id, restaurant_id, name, display_order, created_at
		FROM menu_categories
		WHERE id = $1
	`

	category := &model.MenuCategory{}

	err := db.DB.QueryRow(context.Background(), query, id).Scan(
		&category.ID,
		&category.RestaurantID,
		&category.Name,
		&category.DisplayOrder,
		&category.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func DeleteCategory(id string) error {
	query := `DELETE FROM menu_categories WHERE id = $1`
	_, err := db.DB.Exec(context.Background(), query, id)
	return err
}

// ====== MENU ITEMS ======

func CreateMenuItem(item *model.MenuItem) error {
	query := `
		INSERT INTO menu_items (category_id, name, description, price, image_url, is_veg)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	return db.DB.QueryRow(
		context.Background(),
		query,
		item.CategoryID,
		item.Name,
		item.Description,
		item.Price,
		item.ImageURL,
		item.IsVeg,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
}

func GetMenuItemsByCategory(categoryID string) ([]model.MenuItem, error) {
	query := `
		SELECT id, category_id, name, description, price, image_url, is_available, is_veg, created_at, updated_at
		FROM menu_items
		WHERE category_id = $1
		ORDER BY created_at ASC
	`

	rows, err := db.DB.Query(context.Background(), query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.MenuItem

	for rows.Next() {
		var item model.MenuItem
		err := rows.Scan(
			&item.ID,
			&item.CategoryID,
			&item.Name,
			&item.Description,
			&item.Price,
			&item.ImageURL,
			&item.IsAvailable,
			&item.IsVeg,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func GetMenuItemByID(id string) (*model.MenuItem, error) {
	query := `
		SELECT id, category_id, name, description, price, image_url, is_available, is_veg, created_at, updated_at
		FROM menu_items
		WHERE id = $1
	`

	item := &model.MenuItem{}

	err := db.DB.QueryRow(context.Background(), query, id).Scan(
		&item.ID,
		&item.CategoryID,
		&item.Name,
		&item.Description,
		&item.Price,
		&item.ImageURL,
		&item.IsAvailable,
		&item.IsVeg,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func UpdateMenuItem(id string, req *model.UpdateMenuItemRequest) error {
	query := `
		UPDATE menu_items
		SET name = $1, description = $2, price = $3, image_url = $4, is_available = $5, is_veg = $6, updated_at = NOW()
		WHERE id = $7
	`

	_, err := db.DB.Exec(
		context.Background(),
		query,
		req.Name,
		req.Description,
		req.Price,
		req.ImageURL,
		req.IsAvailable,
		req.IsVeg,
		id,
	)

	return err
}

func DeleteMenuItem(id string) error {
	query := `DELETE FROM menu_items WHERE id = $1`
	_, err := db.DB.Exec(context.Background(), query, id)
	return err
}
