package main

import (
	"fmt"
	"log"
	"net/http"
	"vigilant-spork/db"
	"vigilant-spork/handlers"
	"vigilant-spork/repository"
	"vigilant-spork/routes"
	"vigilant-spork/services"
)

func main() {

	Db := db.InitDb()

	userRepo := &repository.UserRepo{Db: Db}
	productRepo := &repository.ProductRepo{Db: Db}
	cartRepo := &repository.CartRepo{Db: Db}
	orderRepo := &repository.OrderRepo{Db: Db}
	reviewRepo := &repository.ReviewRepo{Db: Db}

	userService := &services.UserService{UserRepo: userRepo}
	productService := &services.ProductService{ProductRepo: productRepo}
	cartService := &services.CartService{CartRepo: cartRepo}
	orderService := &services.OrderService{OrderRepo: orderRepo}
	reviewService := services.NewReviewService(reviewRepo, productRepo)

	userHandler := &handlers.UserHandler{Service: userService}
	productHandler := &handlers.ProductHandler{Service: productService}
	cartHandler := &handlers.CartHandler{Service: cartService}
	orderHandler := &handlers.OrderHandler{Service: orderService}
	reviewHandler := &handlers.ReviewHandler{Service: reviewService}

	r := routes.SetupRouter(userHandler, productHandler, cartHandler, orderHandler, reviewHandler)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("failed to start server", err)
	}
	fmt.Println("server started!")
}