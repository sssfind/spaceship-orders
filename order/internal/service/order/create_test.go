package order_test

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"order/internal/model"
)

// Успешное создание заказа
func (s *OrderServiceTestSuite) TestCreateOrder_Success() {
	userUUID := uuid.New()
	partUUID := uuid.New()

	mockParts := []model.Part{
		{
			UUID:  partUUID,
			Name:  "Тестовое крыло",
			Price: 150.0,
		},
	}

	s.inventoryMock.On("ListParts", mock.Anything, []uuid.UUID{partUUID}).Return(mockParts, nil).Once()

	s.repoMock.On("Create", mock.Anything, mock.AnythingOfType("*model.Order")).Return(nil).Once()

	order, err := s.service.CreateOrder(s.ctx, userUUID, []uuid.UUID{partUUID})

	s.NoError(err)
	s.NotNil(order)
	s.Equal(userUUID, order.UserUUID)
	s.Equal(model.StatusPendingPayment, order.Status)
}

// Склад лежит или вернул ошибку
func (s *OrderServiceTestSuite) TestCreateOrder_InventoryError() {
	userUUID := uuid.New()
	partUUID := uuid.New()

	s.inventoryMock.On("ListParts", mock.Anything, []uuid.UUID{partUUID}).
		Return(nil, errors.New("inventory service unavailable")).Once()

	_, err := s.service.CreateOrder(s.ctx, userUUID, []uuid.UUID{partUUID})

	s.ErrorContains(err, "inventory service unavailable")
}
