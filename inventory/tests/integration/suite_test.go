//go:build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestInventoryIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(InventoryTestSuite))
}
