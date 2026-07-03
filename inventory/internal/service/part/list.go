package part

import (
	"context"
	"inventory/internal/model"
)

func (s *srv) ListParts(ctx context.Context, filter *model.PartsFilter) ([]model.Part, error) {
	return s.partRepo.List(ctx, filter)
}
