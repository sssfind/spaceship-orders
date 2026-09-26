package session

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"iam/internal/model"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, session *model.Session) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO sessions (session_uuid, user_uuid, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
	`, session.SessionUUID, session.UserUUID, session.ExpiresAt, session.CreatedAt)
	return err
}

func (r *Repository) GetByUUID(ctx context.Context, sessionUUID uuid.UUID) (*model.Session, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT session_uuid, user_uuid, expires_at, created_at
		FROM sessions WHERE session_uuid = $1
	`, sessionUUID)

	var session model.Session
	err := row.Scan(&session.SessionUUID, &session.UserUUID, &session.ExpiresAt, &session.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrSessionNotFound
	}
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *Repository) Delete(ctx context.Context, sessionUUID uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM sessions WHERE session_uuid = $1`, sessionUUID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrSessionNotFound
	}
	return nil
}

func (r *Repository) DeleteExpired(ctx context.Context, before time.Time) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM sessions WHERE expires_at < $1`, before)
	return err
}
