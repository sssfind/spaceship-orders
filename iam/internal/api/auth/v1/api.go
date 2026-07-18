package v1

import (
	"iam/internal/service"
	authv1 "spaceship-orders/shared/pkg/proto/auth/v1"
)

type Implementation struct {
	authv1.UnimplementedAuthServiceServer
	authService service.AuthService
}

func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{
		authService: authService,
	}
}
