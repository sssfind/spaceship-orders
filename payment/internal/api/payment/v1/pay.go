package v1

import (
	"context"

	"payment/internal/converter"
	pb "spaceship-orders/shared/pkg/proto/payment/v1"
)

func (h *api) PayOrder(ctx context.Context, req *pb.PayOrderRequest) (*pb.PayOrderResponse, error) {
	domainMethod := converter.ToDomainMethod(req.PaymentMethod)

	txUUID, err := h.paymentService.ProcessPayment(ctx, req.OrderUuid, req.UserUuid, domainMethod)
	if err != nil {
		return nil, err
	}

	return &pb.PayOrderResponse{
		TransactionUuid: txUUID,
	}, nil
}
