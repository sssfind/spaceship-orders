package order

import (
	"context"
	"order/internal/model"

	"github.com/google/uuid"
)

func (s *srv) CancelOrder(ctx context.Context, orderUUID uuid.UUID) error {
	order, err := s.orderRepo.Get(ctx, orderUUID.String())
	if err != nil {
		return err
	}
	if order == nil {
		return model.ErrOrderNotFound
	}

	if order.Status == model.StatusPaid {
		return model.ErrOrderAlreadyPaid
	}

	return s.orderRepo.UpdateStatus(ctx, orderUUID.String(), model.StatusCancelled, "", "")
}
