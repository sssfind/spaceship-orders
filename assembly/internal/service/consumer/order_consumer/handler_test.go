package order_consumer

import (
	"context"
	"platform/pkg/kafka/consumer"
	pbEvents "spaceship-orders/shared/pkg/proto/events/v1"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
)

func (s *HandlerTestSuite) TestHandle_Success() {
	orderID := uuid.New().String()
	userID := uuid.New().String()

	event := &pbEvents.OrderPaidEvent{
		OrderUuid: orderID,
		UserUuid:  userID,
	}
	data, _ := proto.Marshal(event)

	done := make(chan struct{})

	s.producerMock.On("PublishShipAssembled",
		mock.Anything, orderID, userID, mock.AnythingOfType("int64"),
	).Run(func(args mock.Arguments) {
		close(done)
	}).Return(nil).Once()

	msg := consumer.Message{Value: data}

	err := s.handler.Handle(context.Background(), msg)
	s.NoError(err)

	select {
	case <-done:
	case <-time.After(15 * time.Second):
		s.Fail("Timed out waiting for PublishShipAssembled")
	}

	s.producerMock.AssertExpectations(s.T())
}

func (s *HandlerTestSuite) TestHandle_UnmarshalError() {
	msg := consumer.Message{Value: []byte("invalid data")}

	err := s.handler.Handle(context.Background(), msg)

	s.Error(err)
	s.Contains(err.Error(), "failed to unmarshal")
	s.producerMock.AssertNotCalled(s.T(), "PublishShipAssembled", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}
