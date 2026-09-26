//go:build integration

package integration

import (
	"github.com/google/uuid"

	"order/internal/model"
)

func (s *OrderTestSuite) TestDelete_Success() {
	orderUUID := uuid.New()
	err := s.repo.Create(s.ctx, &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   uuid.New(),
		PartUUIDs:  []uuid.UUID{uuid.New()},
		TotalPrice: 100,
		Status:     model.StatusPendingPayment,
	})
	s.Require().NoError(err)

	err = s.repo.Delete(s.ctx, orderUUID.String())
	s.Require().NoError(err)

	got, err := s.repo.Get(s.ctx, orderUUID.String())
	s.ErrorIs(err, model.ErrOrderNotFound)
	s.Nil(got)
}

func (s *OrderTestSuite) TestDelete_NotFound() {
	err := s.repo.Delete(s.ctx, uuid.New().String())
	s.ErrorIs(err, model.ErrOrderNotFound)
}
