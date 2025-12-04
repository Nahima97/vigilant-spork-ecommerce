package services

import (
	"vigilant-spork/models"
	"vigilant-spork/repository"
)

type OrderService struct {
    OrderRepo repository.OrderRepository
}

func (s *OrderService) GetOrderHistory(userID int) ([]models.Order, error) {
	return s.OrderRepo.GetOrderHistory(userID)
}

