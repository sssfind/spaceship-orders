package order

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"order/internal/model"
)

func (s *ServiceSuite) TestListOrders_Success() {
	userUUID := uuid.New()
	want := []*model.Order{
		{OrderUUID: uuid.New(), UserUUID: userUUID, Status: model.StatusPendingPayment},
		{OrderUUID: uuid.New(), UserUUID: userUUID, Status: model.StatusPaid},
	}

	s.repo.EXPECT().
		List(mock.Anything, userUUID).
		Return(want, nil).
		Once()

	got, err := s.svc.ListOrders(s.ctx, userUUID)
	s.Require().NoError(err)
	s.Equal(want, got)
}

func (s *ServiceSuite) TestListOrders_Empty() {
	userUUID := uuid.New()

	s.repo.EXPECT().
		List(mock.Anything, userUUID).
		Return([]*model.Order{}, nil).
		Once()

	got, err := s.svc.ListOrders(s.ctx, userUUID)
	s.Require().NoError(err)
	s.Empty(got)
}
