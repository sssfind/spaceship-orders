package v1

import (
	"context"
	"errors"
	"iam/internal/converter"
	"iam/internal/model"
	authv1 "spaceship-orders/shared/pkg/proto/auth/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Whoami(ctx context.Context, req *authv1.WhoamiRequest) (*authv1.WhoamiResponse, error) {
	user, err := i.authService.Whoami(ctx, req.GetSessionUuid())
	if err != nil {
		if errors.Is(err, model.ErrSessionNotFound) {
			return nil, status.Error(codes.Unauthenticated, "session expired or invalid")
		}
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to recognize session")
	}

	return &authv1.WhoamiResponse{
		User: converter.ToProtoUser(user),
	}, nil
}
