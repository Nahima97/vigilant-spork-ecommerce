package routes

import (
	"net/http"
	"vigilant-spork/handlers"
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

	return r
}
