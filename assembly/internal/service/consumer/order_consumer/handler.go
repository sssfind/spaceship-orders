package order_consumer

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	// ИСПРАВЛЕНО: Импортируем продюсер из пакета assembly, а не order
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
	var event pbEvents.OrderPaidEvent
	if err := proto.Unmarshal(msg.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal OrderPaidEvent: %w", err)
	}

	logger.Info(ctx, fmt.Sprintf("Получено событие оплаты заказа %s. Передаем на верфь для сборки...", event.GetOrderUuid()))

	go func(orderUUID, userUUID string) {
		bgCtx := context.Background()

		buildTime := int64(rand.Intn(10) + 1)
		logger.Info(bgCtx, fmt.Sprintf("Начинается сборка корабля для заказа %s. Время сборки: %d сек.", orderUUID, buildTime))

		time.Sleep(time.Duration(buildTime) * time.Second)

		err := h.producer.PublishShipAssembled(bgCtx, orderUUID, userUUID, buildTime)
		if err != nil {
			logger.Error(bgCtx, fmt.Sprintf("Не удалось отправить событие ShipAssembled для заказа %s: %v", orderUUID, err))
			return
		}

		logger.Info(bgCtx, fmt.Sprintf("Корабль для заказа %s успешно собран и готов к отправке!", orderUUID))
	}(event.GetOrderUuid(), event.GetUserUuid())

	return nil
}
