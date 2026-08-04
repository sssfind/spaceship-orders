package payment_test

import (
	"context"
	"platform/pkg/logger"
	"testing"

	"payment/internal/service"
	paymentService "payment/internal/service/payment"

	"github.com/stretchr/testify/suite"
)

type PaymentServiceTestSuite struct {
	suite.Suite
	ctx     context.Context
	service service.PaymentService
}

func (s *PaymentServiceTestSuite) SetupTest() {
	logger.SetNopLogger()

	s.ctx = context.Background()
	s.service = paymentService.NewService()
}

func TestPaymentServiceSuite(t *testing.T) {
	suite.Run(t, new(PaymentServiceTestSuite))
}
