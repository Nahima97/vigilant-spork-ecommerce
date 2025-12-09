package services

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"vigilant-spork/models"
	"vigilant-spork/repository"
)


type OrderService struct {
	OrderRepo repository.OrderRepository
}

func (s *OrderService) MoveCartToOrder(ctx context.Context, userID uuid.UUID) error {
	err := s.OrderRepo.Transaction(ctx, func(txRepo repository.OrderRepository) error {
		cart, err := txRepo.GetCart(ctx, userID)
		if err != nil {
			return err
		}

		if len(cart.Items) == 0 {
			return fmt.Errorf("cannot create order: cart is empty")
		}

		for _, item := range cart.Items {
			err := txRepo.VerifyAndDeductStock(ctx, &item)
			if err != nil {
				return err
			}
		}

		order, err := txRepo.CreateOrder(ctx, userID)
		if err != nil {
			return err
		}

		err = txRepo.MoveCartItemsToOrder(ctx, order.ID, cart.ID)
		if err != nil {
			return err
		}

		items, err := txRepo.GetOrderItems(ctx, order.ID)
		if err != nil {
			return err
		}

		total := int64(0)
		for _, item := range items {
			total += int64(item.Quantity) * item.UnitPrice
		}

		err = txRepo.UpdateOrderTotal(ctx, total, order.ID)
		if err != nil {
			return err
		}

		err = txRepo.ClearCart(ctx, cart.ID)
		if err != nil {
			return err
		}
		return nil

	})
	return err
}

func (s *OrderService) ShipOrder(ctx context.Context, orderID uuid.UUID) error {
	return s.OrderRepo.Transaction(ctx, func(txRepo repository.OrderRepository) error {
		order, err := txRepo.GetOrder(ctx, orderID)
		if err != nil {
			return err
		}
		if order.Status != "PENDING" {
			return fmt.Errorf("cannot ship order in status %s", order.Status)
		}
		order.Status = "SHIPPED"
		return txRepo.UpdateOrder(ctx, order)
	})
}

func (s *OrderService) CancelOrder(ctx context.Context, orderID uuid.UUID) error {
	return s.OrderRepo.Transaction(ctx, func(txRepo repository.OrderRepository) error {
		order, err := txRepo.GetOrder(ctx, orderID)
		if err != nil {
			return err
		}
		if order.Status != "PENDING" {
			return fmt.Errorf("cannot cancel order in status %s", order.Status)
		}
		order.Status = "CANCELLED"
		return txRepo.UpdateOrder(ctx, order)
	})
}


func (s *OrderService) GetOrderHistory(userID uuid.UUID) ([]models.Order, error) {
	return s.OrderRepo.GetOrderHistory(userID)
}

