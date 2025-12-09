package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"vigilant-spork/middleware"
	"vigilant-spork/models"
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

func (h *OrderHandler) GetOrderHistory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	orders, err := h.Service.GetOrderHistory(userID)
	if err != nil {
		http.Error(w, "Unable to fetch order history", http.StatusInternalServerError)
		return
	}

	// Response for Pending, Shipped & Cancelled
	var response []models.OrderResponse
	for _, o := range orders {

		response = append(response, models.OrderResponse{
			ID:        o.ID,
			Total:     o.Total,
			Status:    o.Status,
			CreatedAt: o.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
