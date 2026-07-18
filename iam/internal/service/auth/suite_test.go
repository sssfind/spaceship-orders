package auth

import (
	"context"
	"testing"

	repoMocks "iam/internal/repository/mocks"
	"iam/internal/service"

	"github.com/stretchr/testify/suite"
)

type AuthServiceTestSuite struct {
	suite.Suite
	ctx         context.Context
	userRepo    *repoMocks.MockUserRepository
	sessionRepo *repoMocks.MockSessionRepository
	service     service.AuthService
}

func (s *AuthServiceTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.userRepo = new(repoMocks.MockUserRepository)
	s.sessionRepo = new(repoMocks.MockSessionRepository)

	s.service = NewAuthService(s.userRepo, s.sessionRepo, nil)
}

func TestAuthServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}
