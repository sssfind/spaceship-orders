//go:build integration

package integration

import (
	"spaceship-orders/shared/pkg/proto/inventory/v1"

	"github.com/google/uuid"
)

func (s *InventoryTestSuite) TestGetPart_NotFound() {
	fakeUUID := uuid.New().String()

	_, err := s.client.GetPart(s.ctx, &v1.GetPartRequest{
		Uuid: fakeUUID,
	})

	s.Error(err)
}
