package v1

import (
	"context"
	"order/internal/model"
	pbPayment "spaceship-orders/shared/pkg/proto/payment/v1"
)

func (c *grpcClient) Pay(ctx context.Context, orderUUID, userUUID string, method model.PaymentMethod) (string, error) {
	var grpcMethod pbPayment.PaymentMethod
	switch method {
	case model.MethodCard:
		grpcMethod = pbPayment.PaymentMethod_PAYMENT_METHOD_CARD
	case model.MethodSbp:
		grpcMethod = pbPayment.PaymentMethod_PAYMENT_METHOD_SBP
	case model.MethodCreditCard:
		grpcMethod = pbPayment.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.MethodInvestorMoney:
		grpcMethod = pbPayment.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		grpcMethod = pbPayment.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}

	res, err := c.client.PayOrder(ctx, &pbPayment.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: grpcMethod,
	})
	if err != nil {
		return "", err
	}

	return res.TransactionUuid, nil
}
