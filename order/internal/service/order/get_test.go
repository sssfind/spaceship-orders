package order_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"order/internal/model"
)

func (s *OrderServiceTestSuite) TestGetOrder_Success() {
	partUUID := uuid.New()
	orderUUID := uuid.New()
	userUUID := uuid.New()
	txUUID := uuid.New()
	payMethod := model.MethodSbp

	mockOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUUIDs:       []uuid.UUID{partUUID},
		TotalPrice:      67,
		Status:          model.StatusPaid,
		TransactionUUID: &txUUID,
		PaymentMethod:   &payMethod,
	}

	s.repoMock.On("Get", mock.Anything, orderUUID.String()).Return(&mockOrder, nil).Once()

	order, err := s.service.GetOrderByUUID(s.ctx, orderUUID)

	s.NoError(err)
	s.NotNil(order)
	s.Equal(orderUUID, order.OrderUUID)
}

func (s *OrderServiceTestSuite) TestGetOrder_NotFound() {
	orderUUID := uuid.New()

	s.repoMock.On("Get", mock.Anything, orderUUID.String()).Return(nil, model.ErrOrderNotFound).Once()

	_, err := s.service.GetOrderByUUID(s.ctx, orderUUID)

	s.ErrorIs(err, model.ErrOrderNotFound)
}
