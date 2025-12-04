package services

import (
	"context"
	"vigilant-spork/repository"

	"github.com/gofrs/uuid"
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

        err = txRepo.ClearCart(ctx, cart.ID)
        if err != nil {
            return err 
        }
        return nil
    })
    return err
}

