package service

import (
	"context"

	"iam/internal/model"
)

type AuthService interface {
	Login(ctx context.Context, login, password string) (string, error)
	Whoami(ctx context.Context, sessionUUID string) (*model.User, error)
}

type UserService interface {
	Register(ctx context.Context, user *model.User, rawPassword string) (string, error)
	GetByUUID(ctx context.Context, userUUID string) (*model.User, error)
}
