//go:build integration

package integration

import (
	"github.com/google/uuid"
	"order/internal/model"
)

func (s *OrderTestSuite) TestList_AfterCreate() {
	userUUID := uuid.New()
	otherUserUUID := uuid.New()

	before, err := s.repo.List(s.ctx, userUUID)
	s.Require().NoError(err)

	orderUUID := uuid.New()
	err = s.repo.Create(s.ctx, &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  []uuid.UUID{uuid.New()},
		TotalPrice: 100,
		Status:     model.StatusPendingPayment,
	})
	s.Require().NoError(err)

	// Чужой заказ не должен попасть в список текущего пользователя.
	err = s.repo.Create(s.ctx, &model.Order{
		OrderUUID:  uuid.New(),
		UserUUID:   otherUserUUID,
		PartUUIDs:  []uuid.UUID{uuid.New()},
		TotalPrice: 100,
		Status:     model.StatusPendingPayment,
	})
	s.Require().NoError(err)

	after, err := s.repo.List(s.ctx, userUUID)
	s.Require().NoError(err)
	s.Equal(len(before)+1, len(after))

	found := false
	for _, o := range after {
		s.Equal(userUUID, o.UserUUID)
		if o.OrderUUID == orderUUID {
			found = true
		}
	}
	s.True(found, "created order must appear in List for its user")
}
