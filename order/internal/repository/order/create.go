package order

import (
	"context"
	"order/internal/model"
)

func (r *repo) Create(ctx context.Context, order *model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[order.OrderUUID.String()] = order
	return nil
}
