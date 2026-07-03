package converter

import (
	"github.com/google/uuid"
	"order/internal/model"
	pbInventory "spaceship-orders/shared/pkg/proto/inventory/v1"
)

func ToDomainParts(pbParts []*pbInventory.Part) []model.Part {
	if pbParts == nil {
		return nil
	}

	parts := make([]model.Part, 0, len(pbParts))
	for _, p := range pbParts {
		parsedUUID, err := uuid.Parse(p.Uuid)
		if err != nil {
			continue
		}

		parts = append(parts, model.Part{
			UUID:  parsedUUID,
			Price: p.Price,
		})
	}

	return parts
}
