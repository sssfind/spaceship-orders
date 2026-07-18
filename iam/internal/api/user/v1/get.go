package v1

import (
	"context"
	"errors"
	"iam/internal/converter"
	"iam/internal/model"
	userv1 "spaceship-orders/shared/pkg/proto/user/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	user, err := i.userService.GetByUUID(ctx, req.GetUserUuid())
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to get user profile")
	}

	return &userv1.GetUserResponse{
		User: converter.ToProtoUser(user),
	}, nil
}
