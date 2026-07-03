package order

import (
	"context"
	"order/internal/model"
	"order/internal/repository"

	"github.com/google/uuid"
)

// InventoryClient описывает требования сервиса к работе со складом
type InventoryClient interface {
	ListParts(ctx context.Context, partUUIDs []uuid.UUID) ([]model.Part, error)
}

// PaymentClient описывает требования сервиса к работе с оплатой
type PaymentClient interface {
	Pay(ctx context.Context, orderUUID, userUUID string, method model.PaymentMethod) (string, error)
}

// srv объединяет в себе репозиторий и все внешние gRPC-клиенты,
type srv struct {
	orderRepo       repository.OrderRepository
	inventoryClient InventoryClient
	paymentClient   PaymentClient
}

func NewService(
	orderRepo repository.OrderRepository,
	inventoryClient InventoryClient,
	paymentClient PaymentClient,
) *srv {
	return &srv{
		orderRepo:       orderRepo,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
