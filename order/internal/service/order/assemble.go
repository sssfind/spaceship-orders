package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"order/internal/model"
)

func (s *srv) AssembleOrder(ctx context.Context, orderUUID uuid.UUID) error {
	// получаем заказ из репозитория для валидации статуса
	order, err := s.orderRepo.Get(ctx, orderUUID.String())
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	// собирать можно только оплаченные заказы
	if order.Status != model.StatusPaid {
		return fmt.Errorf("cannot assemble order %s: current status is %s, expected %s",
			orderUUID, order.Status, model.StatusPaid)
	}

	// безопасно разыменовываем указатели для метода репозитория
	var txUUIDStr string
	if order.TransactionUUID != nil {
		txUUIDStr = order.TransactionUUID.String()
	}

	method := model.MethodUnknown
	if order.PaymentMethod != nil {
		method = *order.PaymentMethod
	}

	err = s.orderRepo.UpdateStatus(
		ctx,
		orderUUID.String(),
		model.StatusAssembled,
		txUUIDStr,
		method,
	)
	if err != nil {
		return fmt.Errorf("failed to update order status to ASSEMBLED: %w", err)
	}

	return nil
}
