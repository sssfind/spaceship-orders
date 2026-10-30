package order

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"order/internal/metrics"
	"order/internal/model"
	"platform/pkg/logger"
	"platform/pkg/tracing"
)

func (s *srv) PayOrder(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	ctx, span := tracing.StartSpan(ctx, "PayOrder")
	defer span.End()

	if s.demoFault == "slow" {
		time.Sleep(2 * time.Second)
	}

	order, err := s.orderRepo.Get(ctx, orderUUID.String())
	if err != nil {
		metrics.OrdersPayFailedTotal.WithLabelValues("repo").Inc()
		return uuid.Nil, err
	}
	if order == nil {
		metrics.OrdersPayFailedTotal.WithLabelValues("not_found").Inc()
		return uuid.Nil, model.ErrOrderNotFound
	}
	if order.Status != model.StatusPendingPayment {
		metrics.OrdersPayFailedTotal.WithLabelValues("status").Inc()
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
		metrics.OrdersPayFailedTotal.WithLabelValues("repo").Inc()
		return uuid.Nil, err
	}

	if err := s.orderRepo.UpdateStatus(ctx, order.OrderUUID.String(), model.StatusPaid, paymentUUID.String(), method); err != nil {
		metrics.OrdersPayFailedTotal.WithLabelValues("repo").Inc()
		return uuid.Nil, err
	}

	metrics.OrdersPaidTotal.WithLabelValues(string(method)).Inc()

	logger.Info(ctx, "order paid",
		zap.String("order_uuid", orderUUID.String()),
		zap.String("payment_uuid", paymentUUID.String()),
		zap.String("payment_method", string(method)),
	)

	return paymentUUID, nil
}
