package user

import (
	"context"
	"encoding/json"

	"iam/internal/model"
)

func (r *userRepo) Create(ctx context.Context, u *model.User) (string, error) {
	methodsJSON, err := json.Marshal(u.NotificationMethods)
	if err != nil {
		return "", err
	}

	query := `
		INSERT INTO users (login, password_hash, email, notification_methods)
		VALUES ($1, $2, $3, $4)
		RETURNING uuid`

	var userUUID string
	err = r.db.QueryRowContext(ctx, query, u.Login, u.PasswordHash, u.Email, methodsJSON).Scan(&userUUID)
	if err != nil {
		return "", err
	}

	return userUUID, nil
}
