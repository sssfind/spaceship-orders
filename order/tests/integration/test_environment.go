//go:build integration

package integration

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"order/internal/repository"
)

type OrderTestSuite struct {
	suite.Suite
	ctx         context.Context
	pgContainer testcontainers.Container
	dbPool      *pgxpool.Pool
	repo        repository.OrderRepository
}
