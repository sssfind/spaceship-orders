package payment_test

import (
	"github.com/google/uuid"
	"payment/internal/model"
)

func (s *PaymentServiceTestSuite) TestProcessPayment_Success() {
	orderUUID := uuid.NewString()
	userUUID := uuid.NewString()
	method := model.MethodCard

	txUUIDStr, err := s.service.ProcessPayment(s.ctx, orderUUID, userUUID, method)

	s.NoError(err)
	s.NotEmpty(txUUIDStr)

	_, parseErr := uuid.Parse(txUUIDStr)
	s.NoError(parseErr, "Сервис обязан генерировать корректный UUID для транзакции")
}
