package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"iam/internal/model"
)

func (s *authService) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.userRepo.GetByLogin(ctx, login)
	if err != nil {
		return "", model.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", model.ErrInvalidCredentials
	}

	sessionUUID := uuid.New().String()
	session := &model.Session{
		SessionUUID: sessionUUID,
		UserUUID:    user.UUID,
		CreatedAt:   time.Now(),
	}

	err = s.sessionRepo.Create(ctx, session, s.sessionCfg.TTL())
	if err != nil {
		return "", err
	}

	return sessionUUID, nil
}
