package v1

import (
	"iam/internal/service"
	userv1 "spaceship-orders/shared/pkg/proto/user/v1"
)

type Implementation struct {
	userv1.UnimplementedUserServiceServer
	userService service.UserService
}

func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
