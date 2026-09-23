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

	// Запрашиваем заказ из базы. Переводим UUID в строку для репозитория
	order, err := s.orderRepo.Get(ctx, orderUUID.String())
	if err != nil {
		return uuid.Nil, err // Ошибка самого репозитория (например, гонка данных)
	}

	// Обязательная проверка на существование заказа, чтобы избежать panic
	if order == nil {
		return uuid.Nil, model.ErrOrderNotFound
	}

	// можно ли оплатить заказ в текущем статусе
	if order.Status != model.StatusPendingPayment {
		return uuid.Nil, model.ErrInvalidOrderStatus
	}

	// Имитируем успешную оплату локально
	txUUID := uuid.New()
	txStr := txUUID.String()

	// фиксируем успешную оплату в репозитории
	err = s.orderRepo.UpdateStatus(ctx, order.OrderUUID.String(), model.StatusPaid, txStr, method)
	if err != nil {
		return uuid.Nil, err
	}

	return txUUID, nil
}
