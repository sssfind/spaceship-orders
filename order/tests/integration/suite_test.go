//go:build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestOrderIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(OrderTestSuite))
}
