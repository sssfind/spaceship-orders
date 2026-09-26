//go:build integration

package integration

import (
	"github.com/google/uuid"
	"order/internal/model"
)

func (s *OrderTestSuite) TestCreate_Success() {
	userUUID := uuid.New()
	partUUID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	orderUUID := uuid.New()

	err := s.repo.Create(s.ctx, &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  []uuid.UUID{partUUID},
		TotalPrice: 100,
		Status:     model.StatusPendingPayment,
		Items: []model.OrderItem{
			{OrderUUID: orderUUID, PartUUID: partUUID, Quantity: 1, UnitPrice: 100},
		},
	})
	s.Require().NoError(err)

	got, err := s.repo.Get(s.ctx, orderUUID.String())
	s.Require().NoError(err)
	s.Require().NotNil(got)
	s.Equal(orderUUID, got.OrderUUID)
	s.Equal(userUUID, got.UserUUID)
	s.Equal([]uuid.UUID{partUUID}, got.PartUUIDs)
	s.Equal(100.0, got.TotalPrice)
	s.Equal(model.StatusPendingPayment, got.Status)
	s.Require().Len(got.Items, 1)
	s.Equal(partUUID, got.Items[0].PartUUID)
	s.Equal(1, got.Items[0].Quantity)
}
