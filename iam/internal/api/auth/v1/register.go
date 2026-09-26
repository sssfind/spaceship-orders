package v1

import (
	"context"
	"errors"

	"iam/internal/model"

	iamV1 "spaceship-orders/shared/pkg/openapi/iam/v1"
)

func (a *api) Register(ctx context.Context, req *iamV1.RegisterRequest) (iamV1.RegisterRes, error) {
	user, err := a.authService.Register(ctx, req.Login, req.Password)
	if err != nil {
		if errors.Is(err, model.ErrUserAlreadyExists) {
			return &iamV1.RegisterConflict{Code: 409, Message: "login already taken"}, nil
		}
		return &iamV1.RegisterInternalServerError{Code: 500, Message: err.Error()}, nil
	}

	return &iamV1.UserDto{
		UserUUID: user.UserUUID,
		Login:    user.Login,
	}, nil
}
