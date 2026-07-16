package order_producer

import (
	"context"
	"fmt"

	"platform/pkg/kafka"
	pbEvents "spaceship-orders/shared/pkg/proto/events/v1"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type OrderProducer interface {
	PublishOrderPaid(ctx context.Context, orderUUID, userUUID, paymentMethod, transactionUUID string) error
}

type orderProducer struct {
	prod kafka.Producer
}

func NewOrderProducer(prod kafka.Producer) *orderProducer {
	return &orderProducer{prod: prod}
}

func (p *orderProducer) PublishOrderPaid(ctx context.Context, orderUUID, userUUID, paymentMethod, transactionUUID string) error {
	event := &pbEvents.OrderPaidEvent{
		EventUuid:       uuid.New().String(),
		OrderUuid:       orderUUID,
		UserUuid:        userUUID,
		PaymentMethod:   paymentMethod,
		TransactionUuid: transactionUUID,
	}

	bytes, err := proto.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal OrderPaidEvent: %w", err)
	}

	err = p.prod.Send(ctx, []byte(orderUUID), bytes)
	if err != nil {
		return fmt.Errorf("failed to send order paid event: %w", err)
	}

	return nil
}
