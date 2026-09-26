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
	ListOrders(ctx context.Context, userUUID uuid.UUID) ([]*model.Order, error)
	DeleteOrder(ctx context.Context, orderUUID uuid.UUID) error
	ListPaymentsByOrder(ctx context.Context, orderUUID uuid.UUID) ([]*model.Payment, error)
	GetPaymentByUUID(ctx context.Context, paymentUUID uuid.UUID) (*model.Payment, error)
}

type PartService interface {
	CreatePart(ctx context.Context, name string, price float64, category string, inStock int) (*model.Part, error)
	GetPartByUUID(ctx context.Context, partUUID uuid.UUID) (*model.Part, error)
	ListParts(ctx context.Context) ([]*model.Part, error)
	UpdatePart(ctx context.Context, partUUID uuid.UUID, name string, price float64, category string, inStock int) (*model.Part, error)
	DeletePart(ctx context.Context, partUUID uuid.UUID) error
}
