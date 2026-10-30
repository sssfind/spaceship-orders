package order

import (
	"order/internal/repository"
	"order/internal/service"
)

type srv struct {
	orderRepo   repository.OrderRepository
	partRepo    repository.PartRepository
	paymentRepo repository.PaymentRepository
	demoFault   string
}

func NewService(
	orderRepo repository.OrderRepository,
	partRepo repository.PartRepository,
	paymentRepo repository.PaymentRepository,
	demoFault string,
) service.OrderService {
	return &srv{
		orderRepo:   orderRepo,
		partRepo:    partRepo,
		paymentRepo: paymentRepo,
		demoFault:   demoFault,
	}
}
