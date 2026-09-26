package order

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"order/internal/model"
)

func (s *ServiceSuite) TestCreateOrder_Success() {
	userUUID := uuid.New()
	part1 := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	part2 := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	partUUIDs := []uuid.UUID{part1, part2}

	s.partRepo.EXPECT().
		GetByUUIDs(mock.Anything, mock.MatchedBy(func(ids []uuid.UUID) bool {
			return len(ids) == 2
		})).
		Return([]*model.Part{
			{UUID: part1, Name: "Ion Thruster", Price: 100, InStock: 10},
			{UUID: part2, Name: "Plasma Shield", Price: 250, InStock: 5},
		}, nil).
		Once()

	s.repo.EXPECT().
		Create(mock.Anything, mock.MatchedBy(func(o *model.Order) bool {
			return o.UserUUID == userUUID &&
				len(o.Items) == 2 &&
				o.TotalPrice == 350.0 &&
				o.Status == model.StatusPendingPayment &&
				o.OrderUUID != uuid.Nil
		})).
		Return(nil).
		Once()

	order, err := s.svc.CreateOrder(s.ctx, userUUID, partUUIDs)
	s.Require().NoError(err)
	s.Require().NotNil(order)
	s.Equal(userUUID, order.UserUUID)
	s.Equal(350.0, order.TotalPrice)
	s.Equal(model.StatusPendingPayment, order.Status)
	s.Len(order.Items, 2)
}

func (s *ServiceSuite) TestCreateOrder_PartNotFound() {
	s.partRepo.EXPECT().
		GetByUUIDs(mock.Anything, mock.Anything).
		Return([]*model.Part{}, nil).
		Once()

	order, err := s.svc.CreateOrder(s.ctx, uuid.New(), []uuid.UUID{uuid.New()})
	s.ErrorIs(err, model.ErrPartNotFound)
	s.Nil(order)
}

func (s *ServiceSuite) TestCreateOrder_RepositoryError() {
	partUUID := uuid.New()
	s.partRepo.EXPECT().
		GetByUUIDs(mock.Anything, mock.Anything).
		Return([]*model.Part{{UUID: partUUID, Price: 100, InStock: 5}}, nil).
		Once()
	s.repo.EXPECT().
		Create(mock.Anything, mock.AnythingOfType("*model.Order")).
		Return(errors.New("db error")).
		Once()

	order, err := s.svc.CreateOrder(s.ctx, uuid.New(), []uuid.UUID{partUUID})
	s.Require().Error(err)
	s.Nil(order)
}

func (s *ServiceSuite) TestCreateOrder_EmptyParts() {
	order, err := s.svc.CreateOrder(s.ctx, uuid.New(), nil)
	s.ErrorIs(err, model.ErrEmptyPartsList)
	s.Nil(order)
}
