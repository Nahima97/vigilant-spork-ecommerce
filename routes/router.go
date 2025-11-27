package routes

import (
	"net/http"
	"os"
	"vigilant-spork/handlers"
	"vigilant-spork/middleware"

	"github.com/gorilla/mux"
)

func SetupRouter(
	userHandler *handlers.UserHandler, productHandler *handlers.ProductHandler, cartHandler *handlers.CartHandler,
	orderHandler *handlers.OrderHandler, reviewHandler *handlers.ReviewHandler) *mux.Router {

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to FutureMarket by Vigilant-Spork!"))
	})

	// User Routes
	r.HandleFunc("/api/v1/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/api/v1/register", userHandler.Register).Methods("POST")

	// Public route: anyone can view reviews
	r.HandleFunc("/products/{productID}/reviews", reviewHandler.GetReviews).Methods("GET")

	// protected routes
	secret := os.Getenv("JWT_SECRET")
	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(middleware.AuthMiddleware(secret))

	// Create a review
	protected.HandleFunc("/products/{product_id}/review", reviewHandler.SubmitReview).Methods("POST")

	// Update a review (authenticated user updates their review)
	protected.HandleFunc("/products/{product_id}/review/{review_id}", reviewHandler.UpdateReview).Methods("PUT")

	// Delete a review (authenticated user deletes their review)
	protected.HandleFunc("/products/{product_id}/review/{review_id}", reviewHandler.DeleteReview).Methods("DELETE")

	// helpful NotFound handler
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - route not found: " + r.URL.Path))
	})

	return r
}
