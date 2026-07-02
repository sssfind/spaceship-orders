package v1

import (
	"context"
	"order/internal/client/converter"
	"order/internal/model"

	pbInventory "spaceship-orders/shared/pkg/proto/inventory/v1"

	"github.com/google/uuid"
)

func (c *grpcClient) ListParts(ctx context.Context, partUUIDs []uuid.UUID) ([]model.Part, error) {
	partUUIDsStrings := make([]string, 0, len(partUUIDs))
	for _, pUUID := range partUUIDs {
		partUUIDsStrings = append(partUUIDsStrings, pUUID.String())
	}

	inventoryReq := &pbInventory.ListPartsRequest{
		Filter: &pbInventory.PartsFilter{
			Uuids: partUUIDsStrings,
		},
	}

	inventoryRes, err := c.client.ListParts(ctx, inventoryReq)
	if err != nil {
		return nil, err
	}
	
	return converter.ToDomainParts(inventoryRes.Parts), nil
}
