package payment

import "payment/internal/service"

type srv struct{}

func NewService() service.PaymentService {
	return &srv{}
}
