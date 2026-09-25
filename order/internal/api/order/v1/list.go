package v1

import (
	"context"
	"order/internal/converter"
	orderV1 "spaceship-orders/shared/pkg/openapi/order/v1"
)

func (h *api) ListOrders(ctx context.Context) (orderV1.ListOrdersRes, error) {
	orders, err := h.orderService.ListOrders(ctx)
	if err != nil {
		return nil, err
	}

	result := make(orderV1.ListOrdersOKApplicationJSON, 0, len(orders))
	for _, o := range orders {
		result = append(result, *converter.OrderToDto(o))
	}
	return &result, nil
}
