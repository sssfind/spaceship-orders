package order

import (
	"context"

	"github.com/google/uuid"
	"order/internal/metrics"
	"order/internal/model"
)

func (s *srv) CreateOrder(ctx context.Context, userUUID uuid.UUID, partUUIDs []uuid.UUID) (*model.Order, error) {
	// Для ДЗ: упрощаем создание заказа, пока нет таблицы деталей
	totalPrice := 100.0 * float64(len(partUUIDs))

	newOrder := &model.Order{
		OrderUUID:  uuid.New(),
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     model.StatusPendingPayment,
	}

	err := s.orderRepo.Create(ctx, newOrder)
	if err != nil {
		return nil, err
	}

	metrics.OrdersTotal.Inc()
	metrics.OrdersRevenueTotal.Add(totalPrice)

	return newOrder, nil
}
