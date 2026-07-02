package order

import (
	"context"
	"order/internal/model"
	"order/internal/repository"

	"github.com/google/uuid"
)

type InventoryClient interface {
	ListParts(ctx context.Context, partUUIDs []uuid.UUID) ([]model.Part, error)
}

type srv struct {
	orderRepo       repository.OrderRepository
	inventoryClient InventoryClient
}

func NewService(orderRepo repository.OrderRepository, inventoryClient InventoryClient) *srv {
	return &srv{
		orderRepo:       orderRepo,
		inventoryClient: inventoryClient,
	}
}
