package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"iam/internal/model"
	"iam/internal/repository"
)

type Service struct {
	users    repository.UserRepository
	sessions repository.SessionRepository
	ttl      time.Duration
}

func NewService(
	users repository.UserRepository,
	sessions repository.SessionRepository,
	ttl time.Duration,
) *Service {
	return &Service{
		users:    users,
		sessions: sessions,
		ttl:      ttl,
	}
}

func (s *Service) Register(ctx context.Context, login, password string) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	user := &model.User{
		UserUUID:     uuid.New(),
		Login:        login,
		PasswordHash: string(hash),
		CreatedAt:    now,
	}

	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) Login(ctx context.Context, login, password string) (uuid.UUID, *model.User, error) {
	user, err := s.users.GetByLogin(ctx, login)
	if err != nil {
		return uuid.Nil, nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return uuid.Nil, nil, model.ErrInvalidCredentials
	}

	now := time.Now().UTC()
	session := &model.Session{
		SessionUUID: uuid.New(),
		UserUUID:    user.UserUUID,
		ExpiresAt:   now.Add(s.ttl),
		CreatedAt:   now,
	}
	if err := s.sessions.Create(ctx, session); err != nil {
		return uuid.Nil, nil, err
	}

	return session.SessionUUID, user, nil
}

func (s *Service) Logout(ctx context.Context, sessionUUID uuid.UUID) error {
	return s.sessions.Delete(ctx, sessionUUID)
}

func (s *Service) Whoami(ctx context.Context, sessionUUID uuid.UUID) (*model.User, error) {
	session, err := s.sessions.GetByUUID(ctx, sessionUUID)
	if err != nil {
		return nil, err
	}

	if time.Now().UTC().After(session.ExpiresAt) {
		_ = s.sessions.Delete(ctx, sessionUUID)
		return nil, model.ErrSessionExpired
	}

	return s.users.GetByUUID(ctx, session.UserUUID)
}
