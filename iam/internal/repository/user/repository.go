package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"iam/internal/model"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, user *model.User) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO users (user_uuid, login, password_hash, created_at)
		VALUES ($1, $2, $3, $4)
	`, user.UserUUID, user.Login, user.PasswordHash, user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (r *Repository) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT user_uuid, login, password_hash, created_at
		FROM users WHERE login = $1
	`, login)

	var user model.User
	err := row.Scan(&user.UserUUID, &user.Login, &user.PasswordHash, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetByUUID(ctx context.Context, userUUID uuid.UUID) (*model.User, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT user_uuid, login, password_hash, created_at
		FROM users WHERE user_uuid = $1
	`, userUUID)

	var user model.User
	err := row.Scan(&user.UserUUID, &user.Login, &user.PasswordHash, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
