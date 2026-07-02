package v1

import (
	"order/internal/service/order"
	pbPayment "spaceship-orders/shared/pkg/proto/payment/v1"
)

type grpcClient struct {
	client pbPayment.PaymentServiceClient
}

func NewClient(client pbPayment.PaymentServiceClient) order.PaymentClient {
	return &grpcClient{client: client}
}
