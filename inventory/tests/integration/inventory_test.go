//go:build integration

package integration

import (
	"github.com/google/uuid"
	"spaceship-orders/shared/pkg/proto/inventory/v1"
)

func (s *InventoryTestSuite) TestGetPart_NotFound() {
	fakeUUID := uuid.New().String()

	_, err := s.client.GetPart(s.ctx, &v1.GetPartRequest{
		Uuid: fakeUUID,
	})

	s.Error(err)
}
