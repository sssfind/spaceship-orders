package v1

import (
	"order/internal/service/order"
	pbInventory "spaceship-orders/shared/pkg/proto/inventory/v1"
)

type grpcClient struct {
	client pbInventory.InventoryServiceClient
}

func NewClient(client pbInventory.InventoryServiceClient) order.InventoryClient {
	return &grpcClient{client: client}
}
