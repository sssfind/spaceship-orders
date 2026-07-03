package v1

import (
	"context"
	"errors"
	"order/internal/converter"
	"order/internal/model"
	orderV1 "spaceship-orders/shared/pkg/openapi/order/v1"
)

func (h *api) GetOrderByUUID(ctx context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	order, err := h.orderService.GetOrderByUUID(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.GenericError{Code: 404, Message: "Order not found"}, nil
		}
		return nil, err
	}

	return converter.OrderToDto(order), nil
}
