package order_producer

import (
	"context"
	"errors"

	pbEvents "spaceship-orders/shared/pkg/proto/events/v1"

	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
)

func (s *ProducerTestSuite) TestPublishShipAssembled_Success() {
	ctx := context.Background()
	orderUUID := "test-order-uuid-111"
	userUUID := "test-user-uuid-222"
	buildTime := int64(7)

	s.kafkaMock.On("Send", ctx, []byte(orderUUID), mock.MatchedBy(func(bytes []byte) bool {
		var event pbEvents.ShipAssembledEvent
		if err := proto.Unmarshal(bytes, &event); err != nil {
			return false
		}
		return event.OrderUuid == orderUUID &&
			event.UserUuid == userUUID &&
			event.BuildTimeSec == buildTime &&
			event.EventUuid != ""
	})).Return(nil).Once()

	err := s.producer.PublishShipAssembled(ctx, orderUUID, userUUID, buildTime)

	s.NoError(err)
	s.kafkaMock.AssertExpectations(s.T())
}

func (s *ProducerTestSuite) TestPublishShipAssembled_SendError() {
	ctx := context.Background()
	orderUUID := "test-order-uuid-111"
	userUUID := "test-user-uuid-222"
	buildTime := int64(4)

	kafkaErr := errors.New("network timeout: broker unreachable")

	s.kafkaMock.On("Send", ctx, []byte(orderUUID), mock.Anything).
		Return(kafkaErr).Once()

	err := s.producer.PublishShipAssembled(ctx, orderUUID, userUUID, buildTime)

	s.Error(err)
	s.Contains(err.Error(), "failed to send ship assembled event via platform producer")
	s.Contains(err.Error(), kafkaErr.Error())
	s.kafkaMock.AssertExpectations(s.T())
}
