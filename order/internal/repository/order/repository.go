package order

import (
	"sync"

	"order/internal/model"
	"order/internal/repository"
)

type repo struct {
	mu     sync.RWMutex
	orders map[string]*model.Order
}

func NewOrderRepository() repository.OrderRepository {
	return &repo{
		orders: make(map[string]*model.Order),
	}
}
