package order

import (
	"context"
	"fmt"

	"order/internal/model"
)

func (r *repo) Create(ctx context.Context, order *model.Order) error {
	query := `
		INSERT INTO orders (order_uuid, user_uuid, part_uuids, total_price, status, transaction_uuid, payment_method)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(ctx, query,
		order.OrderUUID,
		order.UserUUID,
		order.PartUUIDs,
		order.TotalPrice,
		order.Status,
		order.TransactionUUID,
		order.PaymentMethod,
	)
	if err != nil {
		return fmt.Errorf("repository: create order: %w", err)
	}

	return nil
}
