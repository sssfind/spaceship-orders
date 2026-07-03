package order

import (
	"context"

	"order/internal/model"
)

func (r *repo) Get(ctx context.Context, orderUUID string) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[orderUUID]
	if !ok {
		return nil, model.ErrOrderNotFound
	}

	return order, nil
}
