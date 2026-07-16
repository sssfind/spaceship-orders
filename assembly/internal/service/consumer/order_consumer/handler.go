package order_consumer

import (
	"context"
	"fmt"
	"time"

	"assembly/internal/converter/kafka/decoder"
	"assembly/internal/service/producer/order_producer"
	"platform/pkg/kafka/consumer"
	"platform/pkg/logger"
	pbEvents "spaceship-orders/shared/pkg/proto/events/v1"

	"google.golang.org/protobuf/proto"
)

type OrderPaidHandler struct {
	producer order_producer.AssemblyProducer
}

func NewOrderPaidHandler(producer order_producer.AssemblyProducer) *OrderPaidHandler {
	return &OrderPaidHandler{producer: producer}
}

func (h *OrderPaidHandler) Handle(ctx context.Context, msg consumer.Message) error {
	var protoEvent pbEvents.OrderPaidEvent
	if err := proto.Unmarshal(msg.Value, &protoEvent); err != nil {
		return fmt.Errorf("failed to unmarshal OrderPaidEvent: %w", err)
	}

	event := decoder.ToOrderPaidModel(&protoEvent)
	logger.Info(ctx, fmt.Sprintf("[Assembly] Received OrderPaid event for order %s", event.OrderUUID))

	go func(orderUUID, userUUID string) {
		bgCtx := context.Background()

		buildTime := int64(10)
		logger.Info(bgCtx, fmt.Sprintf("[Assembly] Starting construction for order %s. Wait time: %d sec", orderUUID, buildTime))

		time.Sleep(time.Duration(buildTime) * time.Second)

		err := h.producer.PublishShipAssembled(bgCtx, orderUUID, userUUID, buildTime)
		if err != nil {
			logger.Error(bgCtx, fmt.Sprintf("[Assembly] Failed to publish ShipAssembled event for order %s: %v", orderUUID, err))
			return
		}

		logger.Info(bgCtx, fmt.Sprintf("[Assembly] Ship %s successfully assembled and sent!", orderUUID))
	}(event.OrderUUID, event.UserUUID)

	return nil
}
