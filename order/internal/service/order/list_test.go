package order

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"order/internal/model"
)

func (s *ServiceSuite) TestListOrders_Success() {
	want := []*model.Order{
		{OrderUUID: uuid.New(), Status: model.StatusPendingPayment},
		{OrderUUID: uuid.New(), Status: model.StatusPaid},
	}

	s.repo.EXPECT().
		List(mock.Anything).
		Return(want, nil).
		Once()

	got, err := s.svc.ListOrders(s.ctx)
	s.Require().NoError(err)
	s.Equal(want, got)
}

func (s *ServiceSuite) TestListOrders_Empty() {
	s.repo.EXPECT().
		List(mock.Anything).
		Return([]*model.Order{}, nil).
		Once()

	got, err := s.svc.ListOrders(s.ctx)
	s.Require().NoError(err)
	s.Empty(got)
}
