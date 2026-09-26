package part

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"order/internal/model"
	"order/internal/repository"
)

type repo struct {
	db *pgxpool.Pool
}

func NewPartRepository(db *pgxpool.Pool) repository.PartRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, part *model.Part) error {
	query := `
		INSERT INTO parts (part_uuid, name, price, category, in_stock)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query,
		part.UUID,
		part.Name,
		part.Price,
		part.Category,
		part.InStock,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.ErrPartAlreadyExists
		}
		return fmt.Errorf("repository: create part: %w", err)
	}
	return nil
}

func (r *repo) Get(ctx context.Context, partUUID string) (*model.Part, error) {
	query := `
		SELECT part_uuid, name, price, category, in_stock, created_at
		FROM parts
		WHERE part_uuid = $1
	`
	var part model.Part
	err := r.db.QueryRow(ctx, query, partUUID).Scan(
		&part.UUID,
		&part.Name,
		&part.Price,
		&part.Category,
		&part.InStock,
		&part.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrPartNotFound
		}
		return nil, fmt.Errorf("repository: get part: %w", err)
	}
	return &part, nil
}

func (r *repo) GetByUUIDs(ctx context.Context, partUUIDs []uuid.UUID) ([]*model.Part, error) {
	if len(partUUIDs) == 0 {
		return []*model.Part{}, nil
	}

	query := `
		SELECT part_uuid, name, price, category, in_stock, created_at
		FROM parts
		WHERE part_uuid = ANY($1)
	`
	rows, err := r.db.Query(ctx, query, partUUIDs)
	if err != nil {
		return nil, fmt.Errorf("repository: get parts by uuids: %w", err)
	}
	defer rows.Close()

	var parts []*model.Part
	for rows.Next() {
		var part model.Part
		if err := rows.Scan(
			&part.UUID,
			&part.Name,
			&part.Price,
			&part.Category,
			&part.InStock,
			&part.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("repository: scan part: %w", err)
		}
		parts = append(parts, &part)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: get parts rows: %w", err)
	}
	return parts, nil
}

func (r *repo) List(ctx context.Context) ([]*model.Part, error) {
	query := `
		SELECT part_uuid, name, price, category, in_stock, created_at
		FROM parts
		ORDER BY name
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository: list parts: %w", err)
	}
	defer rows.Close()

	parts := make([]*model.Part, 0)
	for rows.Next() {
		var part model.Part
		if err := rows.Scan(
			&part.UUID,
			&part.Name,
			&part.Price,
			&part.Category,
			&part.InStock,
			&part.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("repository: scan part: %w", err)
		}
		parts = append(parts, &part)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: list parts rows: %w", err)
	}
	return parts, nil
}

func (r *repo) Update(ctx context.Context, part *model.Part) error {
	query := `
		UPDATE parts
		SET name = $1, price = $2, category = $3, in_stock = $4
		WHERE part_uuid = $5
	`
	tag, err := r.db.Exec(ctx, query, part.Name, part.Price, part.Category, part.InStock, part.UUID)
	if err != nil {
		return fmt.Errorf("repository: update part: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return model.ErrPartNotFound
	}
	return nil
}

func (r *repo) Delete(ctx context.Context, partUUID string) error {
	query := `DELETE FROM parts WHERE part_uuid = $1`
	tag, err := r.db.Exec(ctx, query, partUUID)
	if err != nil {
		return fmt.Errorf("repository: delete part: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return model.ErrPartNotFound
	}
	return nil
}

func (r *repo) DecrementStock(ctx context.Context, partUUID string, qty int) error {
	query := `
		UPDATE parts
		SET in_stock = in_stock - $1
		WHERE part_uuid = $2 AND in_stock >= $1
	`
	tag, err := r.db.Exec(ctx, query, qty, partUUID)
	if err != nil {
		return fmt.Errorf("repository: decrement stock: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return model.ErrInsufficientStock
	}
	return nil
}
