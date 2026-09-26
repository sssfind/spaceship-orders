//go:build integration

package integration

import (
	"github.com/google/uuid"

	"order/internal/model"
)

func (s *OrderTestSuite) TestList_AfterCreate() {
	before, err := s.repo.List(s.ctx)
	s.Require().NoError(err)

	orderUUID := uuid.New()
	err = s.repo.Create(s.ctx, &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   uuid.New(),
		PartUUIDs:  []uuid.UUID{uuid.New()},
		TotalPrice: 100,
		Status:     model.StatusPendingPayment,
	})
	s.Require().NoError(err)

	after, err := s.repo.List(s.ctx)
	s.Require().NoError(err)
	s.GreaterOrEqual(len(after), len(before)+1)

	found := false
	for _, o := range after {
		if o.OrderUUID == orderUUID {
			found = true
			break
		}
	}
	s.True(found, "created order must appear in List")
}
