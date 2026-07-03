package part

import (
	"context"

	"inventory/internal/model"
)

func (s *srv) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	return s.partRepo.Get(ctx, uuid)
}
