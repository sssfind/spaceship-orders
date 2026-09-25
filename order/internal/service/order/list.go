package order

import (
	"context"
	"order/internal/model"
)

func (s *srv) ListOrders(ctx context.Context) ([]*model.Order, error) {
	return s.orderRepo.List(ctx)
}
