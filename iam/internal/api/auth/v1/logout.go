package v1

import (
	"context"
	"errors"

	"iam/internal/model"

	iamV1 "spaceship-orders/shared/pkg/openapi/iam/v1"
)

func (a *api) Logout(ctx context.Context, params iamV1.LogoutParams) (iamV1.LogoutRes, error) {
	err := a.authService.Logout(ctx, params.XSessionUUID)
	if err != nil {
		if errors.Is(err, model.ErrSessionNotFound) {
			return &iamV1.LogoutUnauthorized{Code: 401, Message: "invalid session"}, nil
		}
		return &iamV1.LogoutInternalServerError{Code: 500, Message: err.Error()}, nil
	}
	return &iamV1.LogoutNoContent{}, nil
}
