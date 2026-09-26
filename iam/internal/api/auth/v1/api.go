package v1

import "iam/internal/service"

type api struct {
	authService service.AuthService
}

func NewAPI(authService service.AuthService) *api {
	return &api{authService: authService}
}
