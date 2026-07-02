package converter

import (
	"order/internal/model"

	pbInventory "spaceship-orders/shared/pkg/proto/inventory/v1"

	"github.com/google/uuid"
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
