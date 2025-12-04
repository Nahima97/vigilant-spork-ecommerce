package handlers

import (
	"encoding/json"
	"net/http"
	"vigilant-spork/middleware"
	"vigilant-spork/models"
	"vigilant-spork/services"
)

type OrderHandler struct {
	Service *services.OrderService
}

func (h *OrderHandler) GetOrderHistory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	orders, err := h.Service.GetOrderHistory(userID)
	if err != nil {
		http.Error(w, "Unable to fetch order history", http.StatusInternalServerError)
		return
	}

	// Response for Pending, Shipped & Cancelled 
	var response []models.OrderResponse
	for _, o := range orders {

		status := "Unknown"
		switch o.Status {

		case "pending", "PENDING":
			status = "Pending"

		case "shipped", "SHIPPED":
			status = "Shipped"

		case "cancelled", "CANCELLED":
			status = "Cancelled"
		}
        //For Postman JSON for the []models.OrderResponse- do i need this?
		response = append(response, models.OrderResponse{
			ID:        o.ID,
			Status:    status,
			CreatedAt: o.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
