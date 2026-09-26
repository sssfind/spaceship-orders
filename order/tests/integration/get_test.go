//go:build integration

package integration

import (
	"github.com/google/uuid"

	"order/internal/model"
)

func (s *OrderTestSuite) TestGet_Success() {
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New()}

	err := s.repo.Create(s.ctx, &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: 100,
		Status:     model.StatusPendingPayment,
	})
	s.Require().NoError(err)

	got, err := s.repo.Get(s.ctx, orderUUID.String())
	s.Require().NoError(err)
	s.Require().NotNil(got)
	s.Equal(orderUUID, got.OrderUUID)
	s.Equal(userUUID, got.UserUUID)
}

func (s *OrderTestSuite) TestGet_NotFound() {
	got, err := s.repo.Get(s.ctx, uuid.New().String())
	s.ErrorIs(err, model.ErrOrderNotFound)
	s.Nil(got)
}
