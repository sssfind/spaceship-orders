package order_producer

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"platform/pkg/kafka"
	pbEvents "spaceship-orders/shared/pkg/proto/events/v1"
)

type AssemblyProducer interface {
	PublishShipAssembled(ctx context.Context, orderUUID, userUUID string, buildTime int64) error
}

type assemblyProducer struct {
	prod kafka.Producer
}

func NewAssemblyProducer(prod kafka.Producer) AssemblyProducer {
	return &assemblyProducer{prod: prod}
}

func (p *assemblyProducer) PublishShipAssembled(ctx context.Context, orderUUID, userUUID string, buildTime int64) error {
	event := &pbEvents.ShipAssembledEvent{
		EventUuid:    uuid.New().String(),
		OrderUuid:    orderUUID,
		UserUuid:     userUUID,
		BuildTimeSec: buildTime,
	}

	bytes, err := proto.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal ShipAssembledEvent: %w", err)
	}

	err = p.prod.Send(ctx, []byte(orderUUID), bytes)
	if err != nil {
		return fmt.Errorf("failed to send ship assembled event via platform producer: %w", err)
	}

	return nil
}
