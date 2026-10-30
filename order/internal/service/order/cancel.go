package order

import (
	"context"

	"github.com/google/uuid"
	"order/internal/metrics"
	"order/internal/model"
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

	if err := s.orderRepo.UpdateStatus(ctx, orderUUID.String(), model.StatusCancelled, "", ""); err != nil {
		return err
	}

	metrics.OrdersCancelledTotal.Inc()
	return nil
}
