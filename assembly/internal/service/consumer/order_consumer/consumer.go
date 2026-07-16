package order_consumer

import (
	"context"
	"fmt"

	"platform/pkg/kafka"
)

type OrderConsumer struct {
	cons    kafka.Consumer
	handler *OrderPaidHandler
}

func NewOrderConsumer(cons kafka.Consumer, handler *OrderPaidHandler) *OrderConsumer {
	return &OrderConsumer{
		cons:    cons,
		handler: handler,
	}
}

func (c *OrderConsumer) Run(ctx context.Context) error {
	err := c.cons.Consume(ctx, c.handler.Handle)
	if err != nil {
		return fmt.Errorf("failed to run assembly consumer: %w", err)
	}
	return nil
}
