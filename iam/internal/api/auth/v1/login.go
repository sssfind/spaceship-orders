package v1

import (
	"context"
	"errors"

	"iam/internal/model"

	iamV1 "spaceship-orders/shared/pkg/openapi/iam/v1"
)

func (a *api) Login(ctx context.Context, req *iamV1.LoginRequest) (iamV1.LoginRes, error) {
	sessionUUID, user, err := a.authService.Login(ctx, req.Login, req.Password)
	if err != nil {
		if errors.Is(err, model.ErrInvalidCredentials) {
			return &iamV1.LoginUnauthorized{Code: 401, Message: "invalid credentials"}, nil
		}
		return &iamV1.LoginInternalServerError{Code: 500, Message: err.Error()}, nil
	}

	return &iamV1.LoginResponse{
		SessionUUID: sessionUUID,
		User: iamV1.UserDto{
			UserUUID: user.UserUUID,
			Login:    user.Login,
		},
	}, nil
}
