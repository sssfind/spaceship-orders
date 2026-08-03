package user

import (
	"context"
	"database/sql"
	"errors"

	"iam/internal/model"
	"iam/internal/repository/converter"
	repoModel "iam/internal/repository/model"
)

func (r *userRepo) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	query := `SELECT id, uuid, login, password_hash, email, notification_methods, created_at FROM users WHERE login = $1`

	var rm repoModel.UserRepoModel
	err := r.db.QueryRowContext(ctx, query, login).Scan(
		&rm.ID, &rm.UUID, &rm.Login, &rm.PasswordHash, &rm.Email, &rm.NotificationMethods, &rm.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return converter.ToUserFromRepo(&rm)
}

func (r *userRepo) GetByUUID(ctx context.Context, userUUID string) (*model.User, error) {
	query := `SELECT id, uuid, login, password_hash, email, notification_methods, created_at FROM users WHERE uuid = $1`

	var rm repoModel.UserRepoModel
	err := r.db.QueryRowContext(ctx, query, userUUID).Scan(
		&rm.ID, &rm.UUID, &rm.Login, &rm.PasswordHash, &rm.Email, &rm.NotificationMethods, &rm.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return converter.ToUserFromRepo(&rm)
}
