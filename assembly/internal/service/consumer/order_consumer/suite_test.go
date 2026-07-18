package order_consumer

import (
	"platform/pkg/logger"
	"testing"

	"assembly/internal/service/producer/order_producer/mocks"

	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite
	producerMock *mocks.AssemblyProducer
	handler      *OrderPaidHandler
}

func (s *HandlerTestSuite) SetupTest() {
	s.producerMock = new(mocks.AssemblyProducer)
	s.handler = NewOrderPaidHandler(s.producerMock)
}

func (s *HandlerTestSuite) SetupSuite() {

	_ = logger.Init("debug", false)
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
