package order_consumer

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"order/internal/service"
	"platform/pkg/kafka/consumer"
	pbEvents "spaceship-orders/shared/pkg/proto/events/v1"
)

type OrderAssembledHandler struct {
	orderSev service.OrderService
}

func NewOrderAssembledHandler(orderSev service.OrderService) *OrderAssembledHandler {
	return &OrderAssembledHandler{orderSev: orderSev}
}

func (h *OrderAssembledHandler) Handle(ctx context.Context, msg consumer.Message) error {
	var event pbEvents.ShipAssembledEvent
	if err := proto.Unmarshal(msg.Value, &event); err != nil {
		return fmt.Errorf("unmarshal event: %w", err)
	}

	orderUUID, err := uuid.Parse(event.GetOrderUuid())
	if err != nil {
		return fmt.Errorf("invalid order uuid format %s: %w", event.GetOrderUuid(), err)
	}

	err = h.orderSev.AssembleOrder(ctx, orderUUID)
	if err != nil {
		return fmt.Errorf("failed to assemble order %s: %w", orderUUID, err)
	}

	return nil
}
