package routes

import (
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"vigilant-spork/handlers"
	"vigilant-spork/middleware"
	"vigilant-spork/services"
)

func SetupRouter(
	userHandler *handlers.UserHandler, productHandler *handlers.ProductHandler, cartHandler *handlers.CartHandler,
	orderHandler *handlers.OrderHandler, reviewHandler *handlers.ReviewHandler, userService *services.UserService) *mux.Router {

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to FutureMarket by Vigilant-Spork!"))
	})

	// Public Routes
	r.HandleFunc("/api/v1/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/api/v1/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/api/v1/products", productHandler.GetProducts).Methods("GET")
	r.HandleFunc("/api/v1/products/{id}", productHandler.GetProductByID).Methods("GET")
	r.HandleFunc("/api/v1/products/{product_id}/reviews", reviewHandler.GetReviews).Methods("GET")

	// Protected routes
	secret := os.Getenv("JWT_SECRET")
	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(middleware.AuthMiddleware(secret))

	protected.HandleFunc("/products", productHandler.AddProduct).Methods("POST")
	protected.HandleFunc("/products/{id}", productHandler.UpdateProduct).Methods("PATCH")
	protected.HandleFunc("/products/{id}", productHandler.DeleteProduct).Methods("DELETE")
	protected.HandleFunc("/cart/{product_id}", cartHandler.AddToCart).Methods("POST")
	protected.HandleFunc("/cart", cartHandler.ViewCart).Methods("GET")
	protected.HandleFunc("/cart/{product_id}", cartHandler.UpdateItemQuantity).Methods("PATCH")
	protected.HandleFunc("/cart/{product_id}", cartHandler.RemoveItem).Methods("DELETE")
	protected.HandleFunc("/checkout", orderHandler.MoveCartToOrder).Methods("POST")
	protected.HandleFunc("/orders", orderHandler.GetOrderHistory).Methods("GET")
	protected.HandleFunc("/products/{product_id}/reviews", reviewHandler.SubmitReview).Methods("POST")
	protected.HandleFunc("/products/{product_id}/review/{review_id}", reviewHandler.UpdateReview).Methods("PATCH")
	protected.HandleFunc("/products/{product_id}/review/{review_id}", reviewHandler.DeleteReview).Methods("DELETE")
	protected.HandleFunc("/logout", userHandler.Logout).Methods("POST")

	// helpful NotFound handler
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - route not found: " + r.URL.Path))
	})

	return r
}
