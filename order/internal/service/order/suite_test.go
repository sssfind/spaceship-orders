package order_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	clientMocks "order/internal/client/grpc/mocks"
	repoMocks "order/internal/repository/mocks"
	"order/internal/service"
	orderService "order/internal/service/order"
)

type OrderServiceTestSuite struct {
	suite.Suite
	ctx           context.Context
	repoMock      *repoMocks.OrderRepository
	inventoryMock *clientMocks.InventoryClient
	paymentMock   *clientMocks.PaymentClient
	service       service.OrderService
}

func (s *OrderServiceTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.repoMock = repoMocks.NewOrderRepository(s.T())
	s.inventoryMock = clientMocks.NewInventoryClient(s.T())
	s.paymentMock = clientMocks.NewPaymentClient(s.T())

	s.service = orderService.NewService(s.repoMock, s.inventoryMock, s.paymentMock)
}

func TestOrderServiceSuite(t *testing.T) {
	suite.Run(t, new(OrderServiceTestSuite))
}
