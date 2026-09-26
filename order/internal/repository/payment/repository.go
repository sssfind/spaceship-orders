package payment

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"order/internal/model"
	"order/internal/repository"
)

type repo struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) repository.PaymentRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, payment *model.Payment) error {
	query := `
		INSERT INTO payments (payment_uuid, order_uuid, amount, method, status)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query,
		payment.PaymentUUID,
		payment.OrderUUID,
		payment.Amount,
		payment.Method,
		payment.Status,
	)
	if err != nil {
		return fmt.Errorf("repository: create payment: %w", err)
	}
	return nil
}

func (r *repo) Get(ctx context.Context, paymentUUID string) (*model.Payment, error) {
	query := `
		SELECT payment_uuid, order_uuid, amount, method, status, created_at
		FROM payments
		WHERE payment_uuid = $1
	`
	var payment model.Payment
	err := r.db.QueryRow(ctx, query, paymentUUID).Scan(
		&payment.PaymentUUID,
		&payment.OrderUUID,
		&payment.Amount,
		&payment.Method,
		&payment.Status,
		&payment.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrPaymentNotFound
		}
		return nil, fmt.Errorf("repository: get payment: %w", err)
	}
	return &payment, nil
}

func (r *repo) ListByOrder(ctx context.Context, orderUUID string) ([]*model.Payment, error) {
	query := `
		SELECT payment_uuid, order_uuid, amount, method, status, created_at
		FROM payments
		WHERE order_uuid = $1
		ORDER BY created_at
	`
	rows, err := r.db.Query(ctx, query, orderUUID)
	if err != nil {
		return nil, fmt.Errorf("repository: list payments: %w", err)
	}
	defer rows.Close()

	payments := make([]*model.Payment, 0)
	for rows.Next() {
		var payment model.Payment
		if err := rows.Scan(
			&payment.PaymentUUID,
			&payment.OrderUUID,
			&payment.Amount,
			&payment.Method,
			&payment.Status,
			&payment.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("repository: scan payment: %w", err)
		}
		payments = append(payments, &payment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: list payments rows: %w", err)
	}
	return payments, nil
}
