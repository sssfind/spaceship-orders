package order

import (
	"context"
	"order/internal/model"

	"github.com/google/uuid"
)

func (s *srv) GetOrderByUUID(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error) {
	order, err := s.orderRepo.Get(ctx, orderUUID.String())
	if err != nil {
		return nil, err
	}
	return order, nil
}
