package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"order/internal/model"
)

func (r *repo) List(ctx context.Context, userUUID uuid.UUID) ([]*model.Order, error) {
	query := `
		SELECT order_uuid, user_uuid, part_uuids, total_price, status, transaction_uuid, payment_method
		FROM orders
		WHERE user_uuid = $1
		ORDER BY order_uuid
	`

	rows, err := r.db.Query(ctx, query, userUUID)
	if err != nil {
		return nil, fmt.Errorf("repository: list orders: %w", err)
	}
	defer rows.Close()

	var orders []*model.Order
	for rows.Next() {
		var order model.Order
		if err := rows.Scan(
			&order.OrderUUID,
			&order.UserUUID,
			&order.PartUUIDs,
			&order.TotalPrice,
			&order.Status,
			&order.TransactionUUID,
			&order.PaymentMethod,
		); err != nil {
			return nil, fmt.Errorf("repository: list orders: %w", err)
		}
		orders = append(orders, &order)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: list orders: %w", err)
	}

	for _, order := range orders {
		items, partUUIDs, err := r.loadItems(ctx, order.OrderUUID.String())
		if err != nil {
			return nil, err
		}
		order.Items = items
		if len(partUUIDs) > 0 {
			order.PartUUIDs = partUUIDs
		}
	}

	return orders, nil
}
