package v1

import (
	"context"
	"errors"

	"order/internal/model"

	orderV1 "spaceship-orders/shared/pkg/openapi/order/v1"
)

func (h *api) DeleteOrder(ctx context.Context, params orderV1.DeleteOrderParams) (orderV1.DeleteOrderRes, error) {
	err := h.orderService.DeleteOrder(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.DeleteOrderNotFound{
				Code:    404,
				Message: "Order not found",
			}, nil
		}
		return nil, err
	}

	return &orderV1.DeleteOrderNoContent{}, nil
}
