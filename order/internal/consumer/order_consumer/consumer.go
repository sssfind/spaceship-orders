package order_consumer

import (
	"context"
	"fmt"

	"platform/pkg/kafka" //[cite: 2]
)

type OrderConsumer struct {
	cons    kafka.Consumer
	handler *OrderAssembledHandler
}

func NewOrderConsumer(cons kafka.Consumer, handler *OrderAssembledHandler) *OrderConsumer {
	return &OrderConsumer{
		cons:    cons,
		handler: handler,
	}
}

func (c *OrderConsumer) Run(ctx context.Context) error {
	err := c.cons.Consume(ctx, c.handler.Handle)
	if err != nil {
		return fmt.Errorf("run order consumer: %w", err)
	}
	return nil
}
