package v1

import (
	"context"
	"errors"

	"iam/internal/model"

	iamV1 "spaceship-orders/shared/pkg/openapi/iam/v1"
)

func (a *api) Whoami(ctx context.Context, params iamV1.WhoamiParams) (iamV1.WhoamiRes, error) {
	user, err := a.authService.Whoami(ctx, params.XSessionUUID)
	if err != nil {
		if errors.Is(err, model.ErrSessionNotFound) ||
			errors.Is(err, model.ErrSessionExpired) ||
			errors.Is(err, model.ErrInvalidCredentials) {
			return &iamV1.WhoamiUnauthorized{Code: 401, Message: "invalid session"}, nil
		}
		return &iamV1.WhoamiInternalServerError{Code: 500, Message: err.Error()}, nil
	}

	return &iamV1.UserDto{
		UserUUID: user.UserUUID,
		Login:    user.Login,
	}, nil
}
