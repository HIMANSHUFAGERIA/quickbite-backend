package repository

import (
	"context"
	"quickbite/db"
	"quickbite/internal/model"
)

func CreateUser(user *model.User) error {
	query := `
		INSERT INTO users (name, email, password, phone, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	return db.DB.QueryRow(
		context.Background(),
		query,
		user.Name,
		user.Email,
		user.Password,
		user.Phone,
		user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func GetUserByEmail(email string) (*model.User, error) {
	query := `
		SELECT id, name, email, password, phone, role, is_verified, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &model.User{}

	err := db.DB.QueryRow(context.Background(), query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Phone,
		&user.Role,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByID(id string) (*model.User, error) {
	query := `
		SELECT id, name, email, phone, role, is_verified, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &model.User{}

	err := db.DB.QueryRow(context.Background(), query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Role,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
