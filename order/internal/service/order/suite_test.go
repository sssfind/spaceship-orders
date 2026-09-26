package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"order/internal/repository/mocks"
	"order/internal/service"
)

type ServiceSuite struct {
	suite.Suite
	ctx  context.Context
	repo *mocks.MockOrderRepository
	svc  service.OrderService
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = mocks.NewMockOrderRepository(s.T())
	s.svc = NewService(s.repo)
}
