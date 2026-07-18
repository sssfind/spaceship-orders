package order_producer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockPlatformProducer struct {
	mock.Mock
}

func (m *MockPlatformProducer) Send(ctx context.Context, key, value []byte) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

type ProducerTestSuite struct {
	suite.Suite
	kafkaMock *MockPlatformProducer
	producer  AssemblyProducer
}

func (s *ProducerTestSuite) SetupTest() {
	s.kafkaMock = new(MockPlatformProducer)
	s.producer = NewAssemblyProducer(s.kafkaMock)
}

func TestProducerTestSuite(t *testing.T) {
	suite.Run(t, new(ProducerTestSuite))
}
