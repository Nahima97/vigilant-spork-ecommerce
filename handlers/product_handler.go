package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"vigilant-spork/middleware"
	"vigilant-spork/models"
	"vigilant-spork/services"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

type ProductHandler struct {
	Service *services.ProductService
}

type ProductResponse struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       string `json:"price"`
	Stock       int     `json:"stock"`
	Data        string  `json:"data"`
	Rating      float64 `json:"rating"`
	Reviews     []models.Review `json:"reviews"`
}

func (h *ProductHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	var products []models.Product
	err := json.NewDecoder(r.Body).Decode(&products)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	role := middleware.GetUserRole(r.Context())
	if role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	err = h.Service.AddProduct(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	productID := mux.Vars(r)["id"]
	productUUID, err := uuid.FromString(productID)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusNotFound)
		return
	}

	product, err := h.Service.GetProductByID(productUUID)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	response := ProductResponse{
		Name:        product.Name,
		Description: product.Description,
		Price:       fmt.Sprintf("%.2f", float64(product.Price)/100),
		Stock:       product.StockQuantity,
		Data:        product.Data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 20
	}

	minPrice, err := strconv.Atoi(r.URL.Query().Get("min_price"))
	if err != nil || minPrice < 1 {
		minPrice = 0
	}

	maxPrice, err := strconv.Atoi(r.URL.Query().Get("max_price"))
	if err != nil || maxPrice < 1 {
		maxPrice = math.MaxInt
	}

	category := r.URL.Query().Get("category")

	totalItems, err := h.Service.GetTotalItems()
	if err != nil {
		http.Error(w, "unable to get total number of items", http.StatusInternalServerError)
		return
	}

	totalPages := (int(totalItems) + limit - 1) / limit

	if totalPages == 0 {
		totalPages = 1
	}

	if page > totalPages {
		page = totalPages
	}

	rawData, err := h.Service.GetProducts(page, limit, minPrice, maxPrice, category)
	if err != nil {
		http.Error(w, "unable to get products", http.StatusInternalServerError)
	}

	type Metadata struct {
		TotalItems  int64 `json:"total_items"`
		TotalPages  int   `json:"total_pages"`
		CurrentPage int   `json:"current_page"`
	}

	var refinedData []ProductResponse
	for _, p := range rawData {
		refinedData = append(refinedData, ProductResponse{
			Name:        p.Name,
			Description: p.Description,
			Price:       fmt.Sprintf("%.2f", float64(p.Price)/100),
			Stock:       p.StockQuantity,
			Data:        p.Data,
		})
	}
	type Response struct {
		Products []ProductResponse
		Metadata Metadata
	}

	response := Response{
		Products: refinedData,
		Metadata: Metadata{
			TotalItems:  totalItems,
			TotalPages:  totalPages,
			CurrentPage: page,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
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
	productUUID, err := uuid.FromString(productID)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusNotFound)
		return
	}

	updatedProduct, err := h.Service.UpdateProduct(productUUID, &product)
	if err != nil {
		http.Error(w, "unable to update product", http.StatusInternalServerError)
		return
	}

	response := ProductResponse {
		Name: updatedProduct.Name,
		Description: updatedProduct.Description,
		Price: fmt.Sprintf("%.2f", float64(updatedProduct.Price)/100),
		Stock: updatedProduct.StockQuantity,
		Data: updatedProduct.Data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {

	role := middleware.GetUserRole(r.Context())
	if role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	productID := mux.Vars(r)["id"]
	productUUID, err := uuid.FromString(productID)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusNotFound)
		return
	}

	err = h.Service.DeleteProduct(productUUID)
	if err != nil {
		http.Error(w, "unable to delete product", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("product deleted successfully!")
}
