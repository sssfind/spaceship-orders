package payment

import (
	"context"
	"fmt"
	"platform/pkg/logger"

	"payment/internal/model"
	"platform/pkg/tracing"

	"github.com/google/uuid"
)

func (s *srv) ProcessPayment(ctx context.Context, orderUUID, userUUID string, method model.PaymentMethod) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "ProcessPayment")
	defer span.End()

	txUUID := uuid.NewString()

	logger.Info(ctx, fmt.Sprintf("Проведена оплата для заказа %s, метод: %s, TX: %s", orderUUID, method, txUUID))

	return txUUID, nil
}
