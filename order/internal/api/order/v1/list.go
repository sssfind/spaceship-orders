package v1

import (
	"context"
	customMiddleware "order/internal/middleware"

	"order/internal/converter"

	orderV1 "spaceship-orders/shared/pkg/openapi/order/v1"
)

func (h *api) ListOrders(ctx context.Context) (orderV1.ListOrdersRes, error) {

	userUUID, ok := customMiddleware.UserUUIDFromContext(ctx)
	if !ok {
		return &orderV1.ListOrdersBadRequest{Code: 401, Message: "Authentication required"}, nil
	}

	orders, err := h.orderService.ListOrders(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	result := make(orderV1.ListOrdersOKApplicationJSON, 0, len(orders))
	for _, o := range orders {
		result = append(result, *converter.OrderToDto(o))
	}
	return &result, nil
}
