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

type CartItemResponse struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	UnitPrice int64     `json:"unit_price"`
}

type ViewCartResponse struct {
	UserID     uuid.UUID          `json:"user_id"`
	Items      []CartItemResponse `json:"items"`
	GrandTotal int64              `json:"grand_total"`
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

func (h *CartHandler) ViewCart(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	cart, err := h.Service.ViewCart(userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"items": []interface{}{},
			"total": 0,
		})
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type CartItemResponse struct {
		ProductID uuid.UUID `json:"product_id"`
		Name      string    `json:"name"`
		Quantity  int       `json:"quantity"`
		UnitPrice int64     `json:"unit_price"`
	}

	var itemsResp []CartItemResponse
	for _, item := range cart.Items {
		itemsResp = append(itemsResp, CartItemResponse{
			ProductID: item.ProductID,
			Name:      item.Product.Name,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		})
	}

	resp := map[string]interface{}{
		"items": itemsResp,
		"total": cart.Total,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
