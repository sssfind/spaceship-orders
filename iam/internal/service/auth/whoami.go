package auth

import (
	"context"
	"iam/internal/model"
)

func (s *authService) Whoami(ctx context.Context, sessionUUID string) (*model.User, error) {
	session, err := s.sessionRepo.Get(ctx, sessionUUID)
	if err != nil {
		return nil, err // Сюда прозрачно пробросится model.ErrSessionNotFound
	}

	return s.userRepo.GetByUUID(ctx, session.UserUUID)
}
