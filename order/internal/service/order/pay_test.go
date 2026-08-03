package order_test

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"order/internal/model"
)

// Успешная оплата
func (s *OrderServiceTestSuite) TestPayOrder_Success() {
	orderUUID := uuid.New()
	userUUID := uuid.New()
	txStr := uuid.New().String()
	method := model.MethodCard

	existingOrder := &model.Order{
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		Status:    model.StatusPendingPayment,
	}

	s.repoMock.On("Get", mock.Anything, orderUUID.String()).Return(existingOrder, nil).Once()
	s.paymentMock.On("Pay", mock.Anything, orderUUID.String(), userUUID.String(), method).Return(txStr, nil).Once()
	s.repoMock.On("UpdateStatus", mock.Anything, orderUUID.String(), model.StatusPaid, txStr, method).Return(nil).Once()

	s.producerMock.On("PublishOrderPaid",
		mock.Anything,      // ctx
		orderUUID.String(), // orderUUID
		userUUID.String(),  // userUUID
		string(method),     // paymentMethod (приводим к строке, как в коде)
		txStr,              // transactionUUID
	).Return(nil).Once() // Успешная отправка

	resTx, err := s.service.PayOrder(s.ctx, orderUUID, method)

	s.NoError(err)
	s.Equal(txStr, resTx.String())

	s.repoMock.AssertExpectations(s.T())
	s.paymentMock.AssertExpectations(s.T())
	s.producerMock.AssertExpectations(s.T())
}

// Заказ не найден в базе данных
func (s *OrderServiceTestSuite) TestPayOrder_OrderNotFound() {
	orderUUID := uuid.New()
	s.repoMock.On("Get", mock.Anything, orderUUID.String()).Return(nil, nil).Once()

	_, err := s.service.PayOrder(s.ctx, orderUUID, model.MethodCard)

	s.ErrorIs(err, model.ErrOrderNotFound)
}

// Заказ имеет невалидный статус для оплаты
func (s *OrderServiceTestSuite) TestPayOrder_InvalidStatus() {
	orderUUID := uuid.New()
	invalidOrder := &model.Order{
		OrderUUID: orderUUID,
		Status:    model.StatusCancelled, // отмененный заказ нельзя оплатить
	}

	s.repoMock.On("Get", mock.Anything, orderUUID.String()).Return(invalidOrder, nil).Once()

	_, err := s.service.PayOrder(s.ctx, orderUUID, model.MethodCard)

	s.ErrorIs(err, model.ErrInvalidOrderStatus)
}

// Внешний сервис оплаты вернул ошибку
func (s *OrderServiceTestSuite) TestPayOrder_PaymentGatewayError() {
	orderUUID := uuid.New()
	userUUID := uuid.New()
	method := model.MethodCard

	existingOrder := &model.Order{
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		Status:    model.StatusPendingPayment,
	}

	s.repoMock.On("Get", mock.Anything, orderUUID.String()).Return(existingOrder, nil).Once()

	// Имитируем падение шлюза банковской системы
	paymentErr := errors.New("insufficient funds")
	s.paymentMock.On("Pay", mock.Anything, orderUUID.String(), userUUID.String(), method).Return("", paymentErr).Once()

	_, err := s.service.PayOrder(s.ctx, orderUUID, method)

	s.ErrorContains(err, "insufficient funds")
}
