package order

import (
	"context"

	"github.com/google/uuid"
	"order/internal/model"
)

func (s *srv) ListPaymentsByOrder(ctx context.Context, orderUUID uuid.UUID) ([]*model.Payment, error) {
	order, err := s.orderRepo.Get(ctx, orderUUID.String())
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, model.ErrOrderNotFound
	}
	return s.paymentRepo.ListByOrder(ctx, orderUUID.String())
}

func (s *srv) GetPaymentByUUID(ctx context.Context, paymentUUID uuid.UUID) (*model.Payment, error) {
	return s.paymentRepo.Get(ctx, paymentUUID.String())
}
