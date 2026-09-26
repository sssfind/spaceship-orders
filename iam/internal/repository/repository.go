package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"iam/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByLogin(ctx context.Context, login string) (*model.User, error)
	GetByUUID(ctx context.Context, userUUID uuid.UUID) (*model.User, error)
}

type SessionRepository interface {
	Create(ctx context.Context, session *model.Session) error
	GetByUUID(ctx context.Context, sessionUUID uuid.UUID) (*model.Session, error)
	Delete(ctx context.Context, sessionUUID uuid.UUID) error
	DeleteExpired(ctx context.Context, before time.Time) error
}
