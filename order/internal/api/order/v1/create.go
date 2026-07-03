package v1

import (
	"context"

	"github.com/google/uuid"
	orderV1 "spaceship-orders/shared/pkg/openapi/order/v1"
)

func (h *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	userUUID, err := uuid.Parse(req.UserUUID)
	if err != nil {
		return &orderV1.CreateOrderBadRequest{Code: 400, Message: "Invalid user_uuid format"}, nil
	}

	newOrder, err := h.orderService.CreateOrder(ctx, userUUID, req.PartUuids)
	if err != nil {
		// Здесь маппим внутренние ошибки сервиса на HTTP статусы OpenAPI
		return &orderV1.CreateOrderInternalServerError{Code: 500, Message: err.Error()}, nil
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  newOrder.OrderUUID,
		TotalPrice: newOrder.TotalPrice,
	}, nil
}
