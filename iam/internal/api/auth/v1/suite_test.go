// iam/internal/api/auth/v1/suite_test.go
package v1

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"iam/internal/service/mocks"
)

type LoginTestSuite struct {
	suite.Suite
	authServiceMock *mocks.MockAuthService
	api             *Implementation
}

func (s *LoginTestSuite) SetupTest() {
	s.authServiceMock = new(mocks.MockAuthService)

	s.api = &Implementation{
		authService: s.authServiceMock,
	}
}

func TestLoginTestSuite(t *testing.T) {
	suite.Run(t, new(LoginTestSuite))
}
