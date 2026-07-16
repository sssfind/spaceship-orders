package order

import (
	"context" // Не забудь добавить импорт пакета context!
	"fmt"

	"github.com/google/uuid"
	"order/internal/model"
)

func (r *repo) UpdateStatus(ctx context.Context, orderUUID string, status model.OrderStatus, txUUID string, method model.PaymentMethod) error {
	query := `
		UPDATE orders
		SET status = $1,
		    transaction_uuid = $2,
		    payment_method = $3
		WHERE order_uuid = $4
	`

	var targetTxUUID *uuid.UUID
	if txUUID != "" {
		parsed, err := uuid.Parse(txUUID)
		if err != nil {
			return fmt.Errorf("repository: неверный формат txUUID: %w", err)
		}
		targetTxUUID = &parsed
	}
	var targetPayMethod *model.PaymentMethod
	if method != "" {
		targetPayMethod = &method
	}

	cmdTag, err := r.db.Exec(ctx, query,
		status,
		targetTxUUID,
		targetPayMethod,
		orderUUID,
	)
	if err != nil {
		return fmt.Errorf("repository: error updating status: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return model.ErrOrderNotFound
	}

	return nil
}
