package repository

import (
	"context"

	"order/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	Get(ctx context.Context, orderUUID string) (*model.Order, error)
	UpdateStatus(ctx context.Context, orderUUID string, status model.OrderStatus, txUUID string, method model.PaymentMethod) error
}
