package order

import (
	"context"

	"github.com/google/uuid"
	"order/internal/model"
	"order/internal/producer/order_producer"
	"order/internal/repository"
	"order/internal/service"
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
	orderProducer   order_producer.OrderProducer
}

func NewService(
	orderRepo repository.OrderRepository,
	inventoryClient InventoryClient,
	paymentClient PaymentClient,
	orderProducer order_producer.OrderProducer,
) service.OrderService {
	return &srv{
		orderRepo:       orderRepo,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		orderProducer:   orderProducer,
	}
}
