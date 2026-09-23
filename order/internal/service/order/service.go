package order

import (
	"order/internal/repository"
	"order/internal/service"
)

// srv объединяет в себе репозиторий
type srv struct {
	orderRepo repository.OrderRepository
}

func NewService(
	orderRepo repository.OrderRepository,
) service.OrderService {
	return &srv{
		orderRepo: orderRepo,
	}
}
