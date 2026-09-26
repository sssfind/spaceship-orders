package part

import (
	"context"

	"github.com/google/uuid"
	"order/internal/model"
	"order/internal/repository"
	"order/internal/service"
)

type srv struct {
	partRepo repository.PartRepository
}

func NewService(partRepo repository.PartRepository) service.PartService {
	return &srv{partRepo: partRepo}
}

func (s *srv) CreatePart(ctx context.Context, name string, price float64, category string, inStock int) (*model.Part, error) {
	if category == "" {
		category = "GENERAL"
	}
	part := &model.Part{
		UUID:     uuid.New(),
		Name:     name,
		Price:    price,
		Category: category,
		InStock:  inStock,
	}
	if err := s.partRepo.Create(ctx, part); err != nil {
		return nil, err
	}
	return part, nil
}

func (s *srv) GetPartByUUID(ctx context.Context, partUUID uuid.UUID) (*model.Part, error) {
	return s.partRepo.Get(ctx, partUUID.String())
}

func (s *srv) ListParts(ctx context.Context) ([]*model.Part, error) {
	return s.partRepo.List(ctx)
}

func (s *srv) UpdatePart(ctx context.Context, partUUID uuid.UUID, name string, price float64, category string, inStock int) (*model.Part, error) {
	part := &model.Part{
		UUID:     partUUID,
		Name:     name,
		Price:    price,
		Category: category,
		InStock:  inStock,
	}
	if err := s.partRepo.Update(ctx, part); err != nil {
		return nil, err
	}
	return s.partRepo.Get(ctx, partUUID.String())
}

func (s *srv) DeletePart(ctx context.Context, partUUID uuid.UUID) error {
	return s.partRepo.Delete(ctx, partUUID.String())
}
