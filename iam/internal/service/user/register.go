package user

import (
	"context"
	"iam/internal/model"

	"golang.org/x/crypto/bcrypt"
)

func (s *userService) Register(ctx context.Context, u *model.User, rawPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	u.PasswordHash = string(hashedPassword)

	return s.userRepo.Create(ctx, u)
}
