package order

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"order/internal/repository"
)

type repo struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) repository.OrderRepository {
	return &repo{
		db: db,
	}
}
