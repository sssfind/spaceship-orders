package order

import (
	"order/internal/repository"
	"order/internal/service"
)

type srv struct {
	orderRepo   repository.OrderRepository
	partRepo    repository.PartRepository
	paymentRepo repository.PaymentRepository
}

func NewService(
	orderRepo repository.OrderRepository,
	partRepo repository.PartRepository,
	paymentRepo repository.PaymentRepository,
) service.OrderService {
	return &srv{
		orderRepo:   orderRepo,
		partRepo:    partRepo,
		paymentRepo: paymentRepo,
	}
}
