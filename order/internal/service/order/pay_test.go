package order

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"order/internal/model"
)

func (s *ServiceSuite) TestPayOrder_Success() {
	orderUUID := uuid.New()
	existing := &model.Order{
		OrderUUID:  orderUUID,
		TotalPrice: 100,
		Status:     model.StatusPendingPayment,
	}

	s.repo.EXPECT().
		Get(mock.Anything, orderUUID.String()).
		Return(existing, nil).
		Once()
	s.paymentRepo.EXPECT().
		Create(mock.Anything, mock.MatchedBy(func(p *model.Payment) bool {
			return p.OrderUUID == orderUUID &&
				p.Amount == 100 &&
				p.Method == model.MethodCard &&
				p.Status == model.PaymentStatusSucceeded
		})).
		Return(nil).
		Once()
	s.repo.EXPECT().
		UpdateStatus(mock.Anything, orderUUID.String(), model.StatusPaid, mock.AnythingOfType("string"), model.MethodCard).
		Return(nil).
		Once()

	txUUID, err := s.svc.PayOrder(s.ctx, orderUUID, model.MethodCard)
	s.Require().NoError(err)
	s.NotEqual(uuid.Nil, txUUID)
}

func (s *ServiceSuite) TestPayOrder_InvalidStatus() {
	orderUUID := uuid.New()
	s.repo.EXPECT().
		Get(mock.Anything, orderUUID.String()).
		Return(&model.Order{OrderUUID: orderUUID, Status: model.StatusPaid}, nil).
		Once()

	txUUID, err := s.svc.PayOrder(s.ctx, orderUUID, model.MethodCard)
	s.ErrorIs(err, model.ErrInvalidOrderStatus)
	s.Equal(uuid.Nil, txUUID)
}

func (s *ServiceSuite) TestPayOrder_NotFound() {
	orderUUID := uuid.New()
	s.repo.EXPECT().
		Get(mock.Anything, orderUUID.String()).
		Return(nil, model.ErrOrderNotFound).
		Once()

	txUUID, err := s.svc.PayOrder(s.ctx, orderUUID, model.MethodCard)
	s.ErrorIs(err, model.ErrOrderNotFound)
	s.Equal(uuid.Nil, txUUID)
}
