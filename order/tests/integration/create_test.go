//go:build integration

package integration

import (
	"github.com/google/uuid"

	"order/internal/model"
)

func (s *OrderTestSuite) TestCreate_Success() {
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New(), uuid.New()}
	orderUUID := uuid.New()

	err := s.repo.Create(s.ctx, &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: 200,
		Status:     model.StatusPendingPayment,
	})
	s.Require().NoError(err)

	got, err := s.repo.Get(s.ctx, orderUUID.String())
	s.Require().NoError(err)
	s.Require().NotNil(got)
	s.Equal(orderUUID, got.OrderUUID)
	s.Equal(userUUID, got.UserUUID)
	s.Equal(partUUIDs, got.PartUUIDs)
	s.Equal(200.0, got.TotalPrice)
	s.Equal(model.StatusPendingPayment, got.Status)
}
