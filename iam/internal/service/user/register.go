package user

import (
	"context"

	"golang.org/x/crypto/bcrypt"
	"iam/internal/model"
)

func (s *userService) Register(ctx context.Context, u *model.User, rawPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	u.PasswordHash = string(hashedPassword)

	return s.userRepo.Create(ctx, u)
}
