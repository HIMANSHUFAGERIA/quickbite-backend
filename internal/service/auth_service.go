package service

import (
	"errors"
	"time"

	"quickbite/config"
	"quickbite/internal/model"
	"quickbite/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(req *model.RegisterRequest, cfg *config.Config) (*model.AuthResponse, error) {
	if req.Role == "" {
		req.Role = "customer"
	}

	if req.Role != "customer" && req.Role != "restaurant_owner" {
		return nil, errors.New("invalid role, must be customer or restaurant_owner")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to process password")
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Phone:    req.Phone,
		Role:     req.Role,
	}

	if err := repository.CreateUser(user); err != nil {
		return nil, errors.New("email already in use")
	}

	token, err := generateJWT(user, cfg)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &model.AuthResponse{Token: token, User: *user}, nil
}

func Login(req *model.LoginRequest, cfg *config.Config) (*model.AuthResponse, error) {
	user, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := generateJWT(user, cfg)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &model.AuthResponse{Token: token, User: *user}, nil
}

func generateJWT(user *model.User, cfg *config.Config) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(cfg.JWTSecret))
}
