package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"order/internal/model"
)

func (r *repo) Create(ctx context.Context, order *model.Order) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("repository: begin create order: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := `
		INSERT INTO orders (order_uuid, user_uuid, part_uuids, total_price, status, transaction_uuid, payment_method)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = tx.Exec(ctx, query,
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

	itemQuery := `
		INSERT INTO order_items (order_uuid, part_uuid, quantity, unit_price)
		VALUES ($1, $2, $3, $4)
	`
	stockQuery := `
		UPDATE parts
		SET in_stock = in_stock - $1
		WHERE part_uuid = $2 AND in_stock >= $1
	`

	for _, item := range order.Items {
		if _, err = tx.Exec(ctx, itemQuery, order.OrderUUID, item.PartUUID, item.Quantity, item.UnitPrice); err != nil {
			return fmt.Errorf("repository: create order item: %w", err)
		}

		tag, err := tx.Exec(ctx, stockQuery, item.Quantity, item.PartUUID)
		if err != nil {
			return fmt.Errorf("repository: decrement stock: %w", err)
		}
		if tag.RowsAffected() == 0 {
			return model.ErrInsufficientStock
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("repository: commit create order: %w", err)
	}
	return nil
}

func (r *repo) loadItems(ctx context.Context, orderUUID string) ([]model.OrderItem, []uuid.UUID, error) {
	query := `
		SELECT order_uuid, part_uuid, quantity, unit_price
		FROM order_items
		WHERE order_uuid = $1
		ORDER BY part_uuid
	`

	rows, err := r.db.Query(ctx, query, orderUUID)
	if err != nil {
		return nil, nil, fmt.Errorf("repository: list order items: %w", err)
	}
	defer rows.Close()

	items := make([]model.OrderItem, 0)
	partUUIDs := make([]uuid.UUID, 0)
	for rows.Next() {
		var item model.OrderItem
		if err := rows.Scan(&item.OrderUUID, &item.PartUUID, &item.Quantity, &item.UnitPrice); err != nil {
			return nil, nil, fmt.Errorf("repository: scan order item: %w", err)
		}
		items = append(items, item)
		for i := 0; i < item.Quantity; i++ {
			partUUIDs = append(partUUIDs, item.PartUUID)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("repository: order items rows: %w", err)
	}
	return items, partUUIDs, nil
}
