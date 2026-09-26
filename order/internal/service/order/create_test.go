package order

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"order/internal/model"
)

func (s *ServiceSuite) TestCreateOrder_Success() {
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New(), uuid.New()}

	s.repo.EXPECT().
		Create(mock.Anything, mock.MatchedBy(func(o *model.Order) bool {
			return o.UserUUID == userUUID &&
				len(o.PartUUIDs) == 2 &&
				o.TotalPrice == 200.0 &&
				o.Status == model.StatusPendingPayment &&
				o.OrderUUID != uuid.Nil
		})).
		Return(nil).
		Once()

	order, err := s.svc.CreateOrder(s.ctx, userUUID, partUUIDs)
	s.Require().NoError(err)
	s.Require().NotNil(order)
	s.Equal(userUUID, order.UserUUID)
	s.Equal(200.0, order.TotalPrice)
	s.Equal(model.StatusPendingPayment, order.Status)
}

func (s *ServiceSuite) TestCreateOrder_RepositoryError() {
	s.repo.EXPECT().
		Create(mock.Anything, mock.AnythingOfType("*model.Order")).
		Return(errors.New("db error")).
		Once()

	order, err := s.svc.CreateOrder(s.ctx, uuid.New(), []uuid.UUID{uuid.New()})
	s.Require().Error(err)
	s.Nil(order)
}
