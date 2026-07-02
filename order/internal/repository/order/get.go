package order

import (
	"order/internal/model"
)

func (r *repo) Get(uuid string) *model.Order {
	r.mu.RLock()
	defer r.mu.RUnlock()
	order, ok := r.orders[uuid]
	if !ok {
		return nil
	}
	return order
}
