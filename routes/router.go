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

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to FutureMarket by Vigilant-Spork!"))
	})

	// User Routes
	r.HandleFunc("/api/v1/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/api/v1/register", userHandler.Register).Methods("POST")

	// protected routes
    secret := os.Getenv("JWT_SECRET")
    protected := r.PathPrefix("/api/v1").Subrouter()
    protected.Use(middleware.AuthMiddleware(secret))

	protected.HandleFunc("/products", productHandler.AddProduct).Methods("POST")
	protected.HandleFunc("/products/{id}", productHandler.UpdateProduct).Methods("PATCH")

	return r
}
