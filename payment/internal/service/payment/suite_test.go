package payment_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"payment/internal/service"
	paymentService "payment/internal/service/payment"
)

type PaymentServiceTestSuite struct {
	suite.Suite
	ctx     context.Context
	service service.PaymentService
}

func (s *PaymentServiceTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.service = paymentService.NewService()
}

func TestPaymentServiceSuite(t *testing.T) {
	suite.Run(t, new(PaymentServiceTestSuite))
}
