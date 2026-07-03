package order

import (
	"context" // Не забудь добавить импорт пакета context!
	"order/internal/model"

	"github.com/google/uuid"
)

func (r *repo) UpdateStatus(ctx context.Context, orderUUID string, status model.OrderStatus, txUUID string, method model.PaymentMethod) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.orders[orderUUID]
	if !ok {
		return nil
	}

	if status != "" {
		order.Status = status
	}

	if txUUID != "" {
		if parsedUUID, err := uuid.Parse(txUUID); err == nil {
			order.TransactionUUID = &parsedUUID
		}
	}

	if method != "" {
		m := method
		order.PaymentMethod = &m
	}
	return nil
}
