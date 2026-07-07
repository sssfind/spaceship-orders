package order

import (
	"order/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) repository.OrderRepository {
	return &repo{
		db: db,
	}
}
