package order

import (
	"context"

	"github.com/google/uuid"
	"order/internal/model"
)

func (s *srv) ListOrders(ctx context.Context, userUUID uuid.UUID) ([]*model.Order, error) {
	return s.orderRepo.List(ctx, userUUID)
}
