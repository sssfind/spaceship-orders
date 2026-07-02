package v1

import (
	"context"
	"errors"
	"order/internal/model"
	orderV1 "spaceship-orders/shared/pkg/openapi/order/v1"
)

func (h *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	var domainMethod model.PaymentMethod
	switch req.PaymentMethod {
	case orderV1.PaymentMethodCARD:
		domainMethod = model.MethodCard
	case orderV1.PaymentMethodSBP:
		domainMethod = model.MethodSbp
	case orderV1.PaymentMethodCREDITCARD:
		domainMethod = model.MethodCreditCard
	case orderV1.PaymentMethodINVESTORMONEY:
		domainMethod = model.MethodInvestorMoney
	default:
		return &orderV1.PayOrderBadRequest{Code: 400, Message: "Unknown payment method"}, nil
	}

	txUUID, err := h.orderService.PayOrder(ctx, params.OrderUUID, domainMethod)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.PayOrderNotFound{Code: 404, Message: "Order not found"}, nil
		}
		if errors.Is(err, model.ErrInvalidOrderStatus) {
			return &orderV1.PayOrderBadRequest{Code: 400, Message: "Invalid order status for payment"}, nil
		}
		return nil, err
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: txUUID,
	}, nil
}
