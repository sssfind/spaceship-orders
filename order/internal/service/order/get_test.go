package order

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"order/internal/model"
)

func (s *ServiceSuite) TestGetOrderByUUID_Success() {
	orderUUID := uuid.New()
	want := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   uuid.New(),
		TotalPrice: 100,
		Status:     model.StatusPendingPayment,
	}

	s.repo.EXPECT().
		Get(mock.Anything, orderUUID.String()).
		Return(want, nil).
		Once()

	got, err := s.svc.GetOrderByUUID(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Equal(want, got)
}

func (s *ServiceSuite) TestGetOrderByUUID_NotFound() {
	orderUUID := uuid.New()
	s.repo.EXPECT().
		Get(mock.Anything, orderUUID.String()).
		Return(nil, model.ErrOrderNotFound).
		Once()

	got, err := s.svc.GetOrderByUUID(s.ctx, orderUUID)
	s.ErrorIs(err, model.ErrOrderNotFound)
	s.Nil(got)
}
