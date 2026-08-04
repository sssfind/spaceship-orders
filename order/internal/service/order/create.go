package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"order/internal/metrics"
	"order/internal/model"
)

func (s *srv) CreateOrder(ctx context.Context, userUUID uuid.UUID, partUUIDs []uuid.UUID) (*model.Order, error) {
	parts, err := s.inventoryClient.ListParts(ctx, partUUIDs)
	if err != nil {
		return nil, err
	}

	if len(parts) != len(partUUIDs) {
		return nil, errors.New("some parts not found")
	}

	var totalPrice float64
	for _, p := range parts {
		totalPrice += p.Price
	}

	newOrder := &model.Order{
		OrderUUID:  uuid.New(),
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     model.StatusPendingPayment,
	}

	err = s.orderRepo.Create(ctx, newOrder)
	if err != nil {
		return nil, err
	}

	metrics.OrdersTotal.Inc()
	metrics.OrdersRevenueTotal.Add(totalPrice)

	return newOrder, nil
}
