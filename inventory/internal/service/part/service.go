package part

import (
	"inventory/internal/repository"
	"inventory/internal/service"
)

type srv struct {
	partRepo repository.PartRepository
}

func NewService(partRepo repository.PartRepository) service.PartService {
	return &srv{
		partRepo: partRepo,
	}
}
