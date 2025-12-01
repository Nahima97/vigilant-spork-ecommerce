package routes

import (
	"net/http"
	"os"
	"vigilant-spork/handlers"
	"vigilant-spork/middleware"
	"vigilant-spork/services"

	"github.com/gorilla/mux"
)

func SetupRouter(
	userHandler *handlers.UserHandler, productHandler *handlers.ProductHandler, cartHandler *handlers.CartHandler,
	orderHandler *handlers.OrderHandler, reviewHandler *handlers.ReviewHandler, userService *services.UserService) *mux.Router {

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to FutureMarket by Vigilant-Spork!"))
	})

	// User Routes
	r.HandleFunc("/api/v1/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/api/v1/register", userHandler.Register).Methods("POST")
	// r.HandleFunc("/api/v1/logout",userHandler.Logout).Methods("POST")
	
	// protected routes
    secret := os.Getenv("JWT_SECRET")
    protected := r.PathPrefix("/api/v1").Subrouter()
   	protected.Use(middleware.AuthMiddleware(secret, userService))

	protected.HandleFunc("/logout",userHandler.Logout).Methods("POST")
   

	return r
}

	
