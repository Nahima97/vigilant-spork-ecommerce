package handlers

import (
	"encoding/json"
	"net/http"
	"vigilant-spork/middleware"
	"vigilant-spork/models"
	"vigilant-spork/services"
)

type ProductHandler struct {
    Service *services.ProductService
}

func (h *ProductHandler) AddProduct (w http.ResponseWriter, r *http.Request) {
    var product models.Product
    err := json.NewDecoder(r.Body).Decode(&product)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return 
    }

    role := middleware.GetUserRole(r.Context())
    if role != "admin" {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }
    
    err = h.Service.AddProduct(&product)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return 
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}