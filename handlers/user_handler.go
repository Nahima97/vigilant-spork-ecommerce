package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"vigilant-spork/middleware"
	"vigilant-spork/models"
	"vigilant-spork/services"
)

type UserHandler struct {
	Service *services.UserService
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var signUp models.User
	err := json.NewDecoder(r.Body).Decode(&signUp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	role := strings.ToLower(signUp.Role)

	if role == "" {
		role = "customer"
	}

	if role != "admin" && role != "customer" {
		http.Error(w, "role must be either 'admin' or 'customer'", http.StatusBadRequest)
		return
	}

	signUp.Role = role

	err = h.Service.RegisterUser(&signUp)
	if err != nil {
		if errors.Is(err, services.ErrEmailExists) {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("registered successfully"))
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var login models.User
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.Service.Login(&login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token := middleware.GetToken(r.Context())
	fmt.Println(token)
	if token == "" {
		http.Error(w, "no token found", http.StatusUnauthorized)
		return
	}

	err := h.Service.UserRepo.AddTokenToBlacklist(token)
	if err != nil {
		http.Error(w, "failed to blacklist token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("logged out successfully"))
}
