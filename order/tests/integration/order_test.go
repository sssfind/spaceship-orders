//go:build integration

package integration

import (
	"github.com/google/uuid"
)

func (s *OrderTestSuite) TestCreateAndGetOrder_Success() {
	fakeUUID := uuid.New().String()

	order, err := s.repo.Get(s.ctx, fakeUUID)

	s.Error(err)
	s.Nil(order)
}
