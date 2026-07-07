package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"order/internal/model"
)

func (r *repo) Get(ctx context.Context, orderUUID string) (*model.Order, error) {
	query := `
		SELECT order_uuid, user_uuid, part_uuids, total_price, status, transaction_uuid, payment_method
		FROM orders
		WHERE order_uuid = $1
	`

	var order model.Order
	err := r.db.QueryRow(ctx, query, orderUUID).Scan(
		&order.OrderUUID,
		&order.UserUUID,
		&order.PartUUIDs,
		&order.TotalPrice,
		&order.Status,
		&order.TransactionUUID,
		&order.PaymentMethod,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrOrderNotFound
		}
		return nil, fmt.Errorf("repository: get order: %w", err)
	}

	return &order, nil
}
