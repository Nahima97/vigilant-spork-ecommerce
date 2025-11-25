package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"vigilant-spork/models"
	"vigilant-spork/services"
)

type UserHandler struct {
    Service *services.UserService
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var signUp models.User
    // Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&signUp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    // Call Service layer
	err := h.Service.RegisterUser(&signUp)
    //Handle errors
	if err != nil {
		if errors.Is(err, services.ErrEmailExists) {
			http.Error(w, err.Error(), http.StatusConflict) // 409 Conflict for duplicate emails
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest) // other validation errors
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(signUp)
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