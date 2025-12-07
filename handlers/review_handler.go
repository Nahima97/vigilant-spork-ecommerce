package handlers

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"vigilant-spork/middleware"
	"vigilant-spork/models"
	"vigilant-spork/services"
)

type ReviewHandler struct {
	Service *services.ReviewService
}

type ReviewResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Rating      int    `json:"rating"`
	Name        string `json:"user_name"`
}

func (h *ReviewHandler) SubmitReview(w http.ResponseWriter, r *http.Request) {
	var review models.Review
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	productID, err := uuid.FromString(vars["product_id"])
	if err != nil {
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}
	review.ProductID = productID

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

	existing.Title = update.Title
	existing.Description = update.Description
	existing.Rating = update.Rating

	if err := h.Service.UpdateReview(existing); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(existing)
}

func (h *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := uuid.FromString(vars["product_id"])
	if err != nil {
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)

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

func (h *ReviewHandler) GetReviews(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := uuid.FromString(vars["product_id"])
	if err != nil {
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	reviews, err := h.Service.GetReviewsForProduct(productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if reviews == nil {
		reviews = []models.Review{}
	}

	var reviewResponses []ReviewResponse
	for _, r := range reviews {
		rr := ReviewResponse{
			Title:       r.Title,
			Description: r.Description,
			Rating:      r.Rating,
			Name:        r.User.Name,
		}
		reviewResponses = append(reviewResponses, rr)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviewResponses)
}
