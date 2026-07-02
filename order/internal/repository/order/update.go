package order

import (
	"order/internal/model"

	"github.com/google/uuid"
)

func (r *repo) UpdateStatus(orderID string, status model.OrderStatus, txUUID string, payMethod model.PaymentMethod) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	order, ok := r.orders[orderID]
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

	if payMethod != "" {
		method := payMethod
		order.PaymentMethod = &method
	}
	return nil
}
