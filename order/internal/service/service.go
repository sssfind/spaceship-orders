package service

import (
	"context"

	"github.com/google/uuid"
	"order/internal/model"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userUUID uuid.UUID, partUUIDs []uuid.UUID) (*model.Order, error)
	GetOrderByUUID(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error)
	CancelOrder(ctx context.Context, orderUUID uuid.UUID) error
	PayOrder(ctx context.Context, orderUUID uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error)
}
