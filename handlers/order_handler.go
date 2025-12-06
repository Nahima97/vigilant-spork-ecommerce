package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"vigilant-spork/middleware"
	"vigilant-spork/repository"
	"vigilant-spork/services"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type OrderHandler struct {
	Service *services.OrderService
}

func (h *OrderHandler) MoveCartToOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetUserID(ctx)
	if userID == uuid.Nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err := h.Service.MoveCartToOrder(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrInsufficientStock):
			http.Error(w, "Insufficient stock", http.StatusConflict)
		case errors.Is(err, gorm.ErrRecordNotFound):
			http.Error(w, "Cart not found", http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Order created successfully")
}
