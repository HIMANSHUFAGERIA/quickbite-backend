package model

import "time"

type User struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Password   string    `json:"-"` // ← "-" means never send password in JSON response
	Phone      string    `json:"phone"`
	Role       string    `json:"role"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// What we receive from the client on register
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Role     string `json:"role"` // customer | restaurant_owner
}

// What we receive from the client on login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// What we send back after successful login/register
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
