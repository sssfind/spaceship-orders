package repository

import (
	"context"

	"github.com/google/uuid"
	"order/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	Get(ctx context.Context, orderUUID string) (*model.Order, error)
	UpdateStatus(ctx context.Context, orderUUID string, status model.OrderStatus, txUUID string, method model.PaymentMethod) error
	List(ctx context.Context, userUUID uuid.UUID) ([]*model.Order, error)
	Delete(ctx context.Context, orderUUID string) error
}

type PartRepository interface {
	Create(ctx context.Context, part *model.Part) error
	Get(ctx context.Context, partUUID string) (*model.Part, error)
	GetByUUIDs(ctx context.Context, partUUIDs []uuid.UUID) ([]*model.Part, error)
	List(ctx context.Context) ([]*model.Part, error)
	Update(ctx context.Context, part *model.Part) error
	Delete(ctx context.Context, partUUID string) error
	DecrementStock(ctx context.Context, partUUID string, qty int) error
}

type PaymentRepository interface {
	Create(ctx context.Context, payment *model.Payment) error
	Get(ctx context.Context, paymentUUID string) (*model.Payment, error)
	ListByOrder(ctx context.Context, orderUUID string) ([]*model.Payment, error)
}
