package repository

import (
	"context"
	"quickbite/db"
	"quickbite/internal/model"
)

func CreateRestaurant(restaurant *model.Restaurant) error {
	query := `
		INSERT INTO restaurants (owner_id, name, description, address, city, image_url)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	return db.DB.QueryRow(
		context.Background(),
		query,
		restaurant.OwnerID,
		restaurant.Name,
		restaurant.Description,
		restaurant.Address,
		restaurant.City,
		restaurant.ImageURL,
	).Scan(&restaurant.ID, &restaurant.CreatedAt, &restaurant.UpdatedAt)
}

func GetRestaurantByID(id string) (*model.Restaurant, error) {
	query := `
		SELECT id, owner_id, name, description, address, city, image_url, is_active, rating, created_at, updated_at
		FROM restaurants
		WHERE id = $1
	`

	restaurant := &model.Restaurant{}

	err := db.DB.QueryRow(context.Background(), query, id).Scan(
		&restaurant.ID,
		&restaurant.OwnerID,
		&restaurant.Name,
		&restaurant.Description,
		&restaurant.Address,
		&restaurant.City,
		&restaurant.ImageURL,
		&restaurant.IsActive,
		&restaurant.Rating,
		&restaurant.CreatedAt,
		&restaurant.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return restaurant, nil
}

func GetRestaurantsByOwner(ownerID string) ([]model.Restaurant, error) {
	query := `
		SELECT id, owner_id, name, description, address, city, image_url, is_active, rating, created_at, updated_at
		FROM restaurants
		WHERE owner_id = $1
		ORDER BY created_at DESC
	`

	rows, err := db.DB.Query(context.Background(), query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var restaurants []model.Restaurant

	for rows.Next() {
		var r model.Restaurant
		err := rows.Scan(
			&r.ID,
			&r.OwnerID,
			&r.Name,
			&r.Description,
			&r.Address,
			&r.City,
			&r.ImageURL,
			&r.IsActive,
			&r.Rating,
			&r.CreatedAt,
			&r.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		restaurants = append(restaurants, r)
	}

	return restaurants, nil
}

func GetAllRestaurants(city string) ([]model.Restaurant, error) {
	var query string
	var rows interface{ Close() }
	var err error

	if city != "" {
		query = `
			SELECT id, owner_id, name, description, address, city, image_url, is_active, rating, created_at, updated_at
			FROM restaurants
			WHERE is_active = true AND city = $1
			ORDER BY rating DESC, created_at DESC
		`
		rows, err = db.DB.Query(context.Background(), query, city)
	} else {
		query = `
			SELECT id, owner_id, name, description, address, city, image_url, is_active, rating, created_at, updated_at
			FROM restaurants
			WHERE is_active = true
			ORDER BY rating DESC, created_at DESC
		`
		rows, err = db.DB.Query(context.Background(), query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.(interface{ Close() }).Close()

	var restaurants []model.Restaurant

	pgxRows := rows.(interface {
		Next() bool
		Scan(dest ...interface{}) error
	})

	for pgxRows.Next() {
		var r model.Restaurant
		err := pgxRows.Scan(
			&r.ID,
			&r.OwnerID,
			&r.Name,
			&r.Description,
			&r.Address,
			&r.City,
			&r.ImageURL,
			&r.IsActive,
			&r.Rating,
			&r.CreatedAt,
			&r.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		restaurants = append(restaurants, r)
	}

	return restaurants, nil
}

func UpdateRestaurant(id string, req *model.UpdateRestaurantRequest) error {
	query := `
		UPDATE restaurants
		SET name = $1, description = $2, address = $3, city = $4, image_url = $5, is_active = $6, updated_at = NOW()
		WHERE id = $7
	`

	_, err := db.DB.Exec(
		context.Background(),
		query,
		req.Name,
		req.Description,
		req.Address,
		req.City,
		req.ImageURL,
		req.IsActive,
		id,
	)

	return err
}

func DeleteRestaurant(id string) error {
	query := `DELETE FROM restaurants WHERE id = $1`
	_, err := db.DB.Exec(context.Background(), query, id)
	return err
}
