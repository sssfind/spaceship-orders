package service

import (
	"context"
	"order/internal/model"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userUUID uuid.UUID, partUUIDs []uuid.UUID) (*model.Order, error)
	GetOrderByUUID(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error)
	CancelOrder(ctx context.Context, orderUUID uuid.UUID) error
	PayOrder(ctx context.Context, orderUUID uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error)
}
