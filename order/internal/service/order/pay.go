package order

import (
	"context"

	"github.com/google/uuid"
	"order/internal/model"
	"platform/pkg/tracing"
)

func (s *srv) PayOrder(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	ctx, span := tracing.StartSpan(ctx, "PayOrder")
	defer span.End()

	order, err := s.orderRepo.Get(ctx, orderUUID.String())
	if err != nil {
		return uuid.Nil, err
	}
	if order == nil {
		return uuid.Nil, model.ErrOrderNotFound
	}
	if order.Status != model.StatusPendingPayment {
		return uuid.Nil, model.ErrInvalidOrderStatus
	}

	paymentUUID := uuid.New()
	payment := &model.Payment{
		PaymentUUID: paymentUUID,
		OrderUUID:   order.OrderUUID,
		Amount:      order.TotalPrice,
		Method:      method,
		Status:      model.PaymentStatusSucceeded,
	}

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return uuid.Nil, err
	}

	if err := s.orderRepo.UpdateStatus(ctx, order.OrderUUID.String(), model.StatusPaid, paymentUUID.String(), method); err != nil {
		return uuid.Nil, err
	}

	return paymentUUID, nil
}
