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

func (i *Implementation) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	u := &model.User{
		Login:               req.GetLogin(),
		Email:               req.GetEmail(),
		NotificationMethods: converter.ToModelNotificationMethods(req.GetNotificationMethods()),
	}

	userUUID, err := i.userService.Register(ctx, u, req.GetPassword())
	if err != nil {
		if errors.Is(err, model.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &userv1.RegisterResponse{
		UserUuid: userUUID,
	}, nil
}
