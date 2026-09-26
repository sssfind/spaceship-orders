package order

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"order/internal/model"
)

func (s *ServiceSuite) TestCancelOrder_Success() {
	orderUUID := uuid.New()
	s.repo.EXPECT().
		Get(mock.Anything, orderUUID.String()).
		Return(&model.Order{OrderUUID: orderUUID, Status: model.StatusPendingPayment}, nil).
		Once()
	s.repo.EXPECT().
		UpdateStatus(mock.Anything, orderUUID.String(), model.StatusCancelled, "", model.PaymentMethod("")).
		Return(nil).
		Once()

	err := s.svc.CancelOrder(s.ctx, orderUUID)
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestCancelOrder_AlreadyPaid() {
	orderUUID := uuid.New()
	s.repo.EXPECT().
		Get(mock.Anything, orderUUID.String()).
		Return(&model.Order{OrderUUID: orderUUID, Status: model.StatusPaid}, nil).
		Once()

	err := s.svc.CancelOrder(s.ctx, orderUUID)
	s.ErrorIs(err, model.ErrOrderAlreadyPaid)
}

func (s *ServiceSuite) TestCancelOrder_NotFound() {
	orderUUID := uuid.New()
	s.repo.EXPECT().
		Get(mock.Anything, orderUUID.String()).
		Return(nil, model.ErrOrderNotFound).
		Once()

	err := s.svc.CancelOrder(s.ctx, orderUUID)
	s.ErrorIs(err, model.ErrOrderNotFound)
}
