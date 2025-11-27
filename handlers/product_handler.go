package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
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

func (h *ProductHandler) GetProductByID (w http.ResponseWriter, r *http.Request) {
    productID := mux.Vars(r)["id"]
    productUUID, err := uuid.Parse(productID)
    if err != nil {
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return 
    }

    product, err := h.Service.GetProductByID(productUUID)
    if err != nil {
        http.Error(w, "Product ID not found", http.StatusNotFound)
        return 
    }

    type ProductResponse struct {
        Name        string    `json:"name"`
        Description string    `json:"description"`
        Price       float64   `json:"price"`
        Stock       int       `json:"stock"`
        Data        string    `json:"data"`
    }

    productresponse := ProductResponse{
        Name: product.Name,
        Description: product.Description,
        Price: product.Price,
        Stock: product.StockQuantity,
        Data: product.Data,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(productresponse)
}

func (h *ProductHandler) GetProducts (w http.ResponseWriter, r *http.Request) {
    page, err := strconv.Atoi(r.URL.Query().Get("page"))
    if err != nil || page < 1 {
        page = 1 // defaults to page 1
    }

    limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
    if err != nil || limit < 1 {
        limit = 20 // default to 20 items per page 
    }

    products, err := h.Service.GetProducts(page, limit)
    if err != nil {
        http.Error(w, "Unable to get products", http.StatusInternalServerError)
        return 
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(products)
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

