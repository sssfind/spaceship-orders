package order

import (
	"context"
	"fmt"

	"order/internal/model"

	"github.com/google/uuid"
)

func (s *srv) PayOrder(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
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

	// Вызываем gRPC-клиент оплаты, передавая нужные для контракта строки
	txStr, err := s.paymentClient.Pay(ctx, order.OrderUUID.String(), order.UserUUID.String(), method)
	if err != nil {
		// Если оплата не прошла - прерываем выполнение, статус заказа остается PENDING
		return uuid.Nil, err
	}

	// Конвертируем string транзакции от gRPC в доменный тип uuid.UUID
	txUUID, err := uuid.Parse(txStr)
	if err != nil {
		return uuid.Nil, err
	}

	// фиксируем успешную оплату в репозитории
	err = s.orderRepo.UpdateStatus(ctx, order.OrderUUID.String(), model.StatusPaid, txStr, method)
	if err != nil {
		return uuid.Nil, err
	}

	err = s.orderProducer.PublishOrderPaid(ctx, order.OrderUUID.String(), order.UserUUID.String(), string(method), txStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to publish order paid event: %w", err)
	}

	return txUUID, nil
}
