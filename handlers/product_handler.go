package handlers

import (
	"encoding/json"
	"net/http"
	"vigilant-spork/middleware"
	"vigilant-spork/models"
	"vigilant-spork/services"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

func (h *ProductHandler) UpdateProduct (w http.ResponseWriter, r *http.Request) {
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

    productID := mux.Vars(r)["id"]
    productUUID, err := uuid.Parse(productID)
    if err != nil {
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return 
    }
    
    updatedProduct, err := h.Service.UpdateProduct(productUUID, &product)
    if err != nil {
        http.Error(w, "unable to update product", http.StatusInternalServerError)
        return 
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedProduct)
}

func (h *ProductHandler) DeleteProduct (w http.ResponseWriter, r *http.Request) {
    
    role := middleware.GetUserRole(r.Context())
    if role != "admin" {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }

    productID := mux.Vars(r)["id"]
    productUUID, err := uuid.Parse(productID)
    if err != nil {
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return 
    }

    err = h.Service.DeleteProduct(productUUID)
    if err != nil {
        http.Error(w, "unable to delete product", http.StatusInternalServerError)
        return 
    }
    w.WriteHeader(http.StatusOK)
}