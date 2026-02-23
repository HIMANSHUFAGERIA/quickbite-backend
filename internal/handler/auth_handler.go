package handler

import (
	"encoding/json"
	"net/http"

	"quickbite/config"
	"quickbite/internal/model"
	"quickbite/internal/service"
	"quickbite/internal/utils"
)

type AuthHandler struct {
	cfg *config.Config
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		utils.WriteError(w, http.StatusBadRequest, "name, email and password are required")
		return
	}

	if len(req.Password) < 6 {
		utils.WriteError(w, http.StatusBadRequest, "password must be at least 6 characters")
		return
	}

	resp, err := service.Register(&req, h.cfg)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		utils.WriteError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	resp, err := service.Login(&req, h.cfg)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}
