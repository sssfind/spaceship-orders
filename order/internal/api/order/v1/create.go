package v1

import (
	"context"
	"errors"

	"order/internal/model"
	customMiddleware "order/internal/middleware"

	orderV1 "spaceship-orders/shared/pkg/openapi/order/v1"
)

func (h *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	userUUID, ok := customMiddleware.UserUUIDFromContext(ctx)
	if !ok {
		return &orderV1.CreateOrderBadRequest{Code: 401, Message: "Authentication required"}, nil
	}

	newOrder, err := h.orderService.CreateOrder(ctx, userUUID, req.PartUuids)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrEmptyPartsList),
			errors.Is(err, model.ErrPartNotFound),
			errors.Is(err, model.ErrInsufficientStock):
			return &orderV1.CreateOrderBadRequest{Code: 400, Message: err.Error()}, nil
		default:
			return &orderV1.CreateOrderInternalServerError{Code: 500, Message: err.Error()}, nil
		}
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  newOrder.OrderUUID,
		TotalPrice: newOrder.TotalPrice,
	}, nil
}
