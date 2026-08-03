package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"iam/internal/model"
	authv1 "spaceship-orders/shared/pkg/proto/auth/v1"
)

func (i *Implementation) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	sessionUUID, err := i.authService.Login(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		if errors.Is(err, model.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &authv1.LoginResponse{
		SessionUuid: sessionUUID,
	}, nil
}
