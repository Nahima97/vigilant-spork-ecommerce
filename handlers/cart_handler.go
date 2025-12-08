package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"vigilant-spork/middleware"
	"vigilant-spork/repository"
	"vigilant-spork/services"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type CartHandler struct {
    Service *services.CartService
}

func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
    productID := mux.Vars(r)["product_id"]
	productUUID, err := uuid.FromString(productID)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

    userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		http.Error(w, "no userID found", http.StatusInternalServerError)
		return
	}

    err = h.Service.AddToCart(userID, productUUID)
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
	json.NewEncoder(w).Encode("Item added to cart successfully")
}