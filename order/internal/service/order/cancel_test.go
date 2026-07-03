package order_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"order/internal/model"
)

// Успешная отмена заказа
func (s *OrderServiceTestSuite) TestCancelOrder_Success() {
	orderUUID := uuid.New()

	existingOrder := &model.Order{
		OrderUUID: orderUUID,
		Status:    model.StatusPendingPayment,
	}

	s.repoMock.On("Get", mock.Anything, orderUUID.String()).Return(existingOrder, nil).Once()
	s.repoMock.On("UpdateStatus", mock.Anything, orderUUID.String(), model.StatusCancelled, "", model.PaymentMethod("")).Return(nil).Once()

	err := s.service.CancelOrder(s.ctx, orderUUID)

	s.NoError(err)
}

// Заказ не найден
func (s *OrderServiceTestSuite) TestCancelOrder_NotFound() {
	orderUUID := uuid.New()
	s.repoMock.On("Get", mock.Anything, orderUUID.String()).Return(nil, nil).Once()

	err := s.service.CancelOrder(s.ctx, orderUUID)

	s.ErrorIs(err, model.ErrOrderNotFound)
}

// Нельзя отменить оплаченный заказ
func (s *OrderServiceTestSuite) TestCancelOrder_AlreadyPaid() {
	orderUUID := uuid.New()

	paidOrder := &model.Order{
		OrderUUID: orderUUID,
		Status:    model.StatusPaid,
	}

	s.repoMock.On("Get", mock.Anything, orderUUID.String()).Return(paidOrder, nil).Once()

	err := s.service.CancelOrder(s.ctx, orderUUID)

	s.ErrorIs(err, model.ErrOrderAlreadyPaid)
}
