//go:build integration

package integration

import (
	"context"
)

func (s *InventoryTestSuite) TearDownSuite() {
	if s.grpcConn != nil {
		_ = s.grpcConn.Close()
	}
	if s.mongoContainer != nil {
		_ = s.mongoContainer.Terminate(context.Background())
	}
}
