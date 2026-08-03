package repository

import (
	"context"
	"time"

	"iam/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (string, error)
	GetByLogin(ctx context.Context, login string) (*model.User, error)
	GetByUUID(ctx context.Context, userUUID string) (*model.User, error)
}

type SessionRepository interface {
	Create(ctx context.Context, session *model.Session, ttl time.Duration) error
	Get(ctx context.Context, sessionUUID string) (*model.Session, error)
}
