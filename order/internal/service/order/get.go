package order

import (
	"context"

	"github.com/google/uuid"
	"order/internal/model"
)

func (s *srv) GetOrderByUUID(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error) {
	order, err := s.orderRepo.Get(ctx, orderUUID.String())
	if err != nil {
		return nil, err
	}
	return order, nil
}
