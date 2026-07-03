package part

import (
	"context"
	"inventory/internal/model"
	"inventory/internal/repository/converter"
)

// Реализация метода Get для интерфейса
func (r *repo) Get(ctx context.Context, uuid string) (*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, exists := r.parts[uuid]
	if !exists {
		return nil, model.ErrPartNotFound
	}

	domainPart := converter.PartToDomain(part)

	return &domainPart, nil
}