package order

import (
	"context"

	"order/internal/model"

	"github.com/google/uuid"
)

func (s *srv) ListOrders(ctx context.Context, userUUID uuid.UUID) ([]*model.Order, error) {
	return s.orderRepo.List(ctx, userUUID)
}
