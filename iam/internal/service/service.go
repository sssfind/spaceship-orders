package service

import (
	"context"

	"github.com/google/uuid"
	"iam/internal/model"
)

type AuthService interface {
	Register(ctx context.Context, login, password string) (*model.User, error)
	Login(ctx context.Context, login, password string) (sessionUUID uuid.UUID, user *model.User, err error)
	Logout(ctx context.Context, sessionUUID uuid.UUID) error
	Whoami(ctx context.Context, sessionUUID uuid.UUID) (*model.User, error)
}
