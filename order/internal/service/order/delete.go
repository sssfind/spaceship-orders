package order

import (
	"context"

	"github.com/google/uuid"
)

func (s *srv) DeleteOrder(ctx context.Context, orderUUID uuid.UUID) error {
	return s.orderRepo.Delete(ctx, orderUUID.String())
}
