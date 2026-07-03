package v1

import (
	"payment/internal/service"
	pb "spaceship-orders/shared/pkg/proto/payment/v1"
)

type api struct {
	pb.UnimplementedPaymentServiceServer
	paymentService service.PaymentService
}

func NewAPI(paymentService service.PaymentService) *api {
	return &api{
		paymentService: paymentService,
	}
}
