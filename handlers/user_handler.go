package handlers

import (
	"encoding/json"
	"net/http"
	"vigilant-spork/models"
	"vigilant-spork/services"
)

type UserHandler struct {
    Service *services.UserService
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
