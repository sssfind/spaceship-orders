package user

import (
	"context"
	"iam/internal/model"
)

func (s *userService) GetByUUID(ctx context.Context, userUUID string) (*model.User, error) {
	return s.userRepo.GetByUUID(ctx, userUUID)
}
