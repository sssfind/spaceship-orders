package auth

import (
	"iam/internal/config"
	"iam/internal/repository"
	"iam/internal/service"
)

type authService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	sessionCfg  config.SessionConfig
}

func NewAuthService(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	sessionCfg config.SessionConfig,
) service.AuthService {
	return &authService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		sessionCfg:  sessionCfg,
	}
}
