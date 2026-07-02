package v1

import (
	"context"
	"errors"
	"order/internal/model"
	orderV1 "spaceship-orders/shared/pkg/openapi/order/v1"
)

func (h *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	err := h.orderService.CancelOrder(ctx, params.OrderUUID)
	if err != nil {

		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.CancelOrderNotFound{
				Code:    404,
				Message: "Order not found",
			}, nil
		}

		if errors.Is(err, model.ErrOrderAlreadyPaid) {
			return &orderV1.CancelOrderConflict{
				Code:    409,
				Message: "Order already paid and cannot be cancelled",
			}, nil
		}

		return nil, err
	}

	return &orderV1.CancelOrderNoContent{}, nil
}
