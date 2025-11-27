package handlers

import (
	"encoding/json"
	"net/http"
	"vigilant-spork/middleware"
	"vigilant-spork/models"
	"vigilant-spork/services"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

type ReviewHandler struct {
	Service *services.ReviewService
}

// SubmitReview handles POST /products/{product_id}/reviews
func (h *ReviewHandler) SubmitReview(w http.ResponseWriter, r *http.Request) {
	var review models.Review
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// URL param product_id
	vars := mux.Vars(r)
	productID, err := uuid.FromString(vars["product_id"])
	if err != nil {
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}
	review.ProductID = productID

	// Get authenticated user ID
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	review.UserID = userID

	if err := h.Service.SubmitReview(&review); err != nil {
		switch err {
		case services.ErrInvalidRating:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case services.ErrRateLimitExceeded:
			http.Error(w, err.Error(), http.StatusTooManyRequests)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(review)
}

// UpdateReview handles PUT /products/{product_id}/reviews
func (h *ReviewHandler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	var update models.Review
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	productID, err := uuid.FromString(vars["product_id"])
	if err != nil {
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)

	// Fetch existing review
	existing, err := h.Service.GetReviewByUserForProduct(userID, productID)
	if err != nil {
		http.Error(w, "review not found", http.StatusNotFound)
		return
	}

	// Update fields
	existing.Title = update.Title
	existing.Description = update.Description
	existing.Rating = update.Rating

	if err := h.Service.UpdateReview(existing); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(existing)
}

// DeleteReview handles DELETE /products/{product_id}/reviews
func (h *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := uuid.FromString(vars["product_id"])
	if err != nil {
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)

	// Fetch existing review
	existing, err := h.Service.GetReviewByUserForProduct(userID, productID)
	if err != nil {
		http.Error(w, "review not found", http.StatusNotFound)
		return
	}

	if err := h.Service.DeleteReview(existing.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetReviews handles GET /products/{product_id}/reviews
func (h *ReviewHandler) GetReviews(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
	// r.URL.Query().Get("product_id")
    productID, err := uuid.FromString(vars["productID"])
    if err != nil {
        http.Error(w, "invalid product ID", http.StatusBadRequest)
        return
    }

    // Fetch all reviews for this product
    reviews, err := h.Service.GetReviewsForProduct(productID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // If no reviews exist, just return an empty slice (not nil)
    if reviews == nil {
        reviews = []models.Review{}
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(reviews)
}

