//go:build integration

package integration

import (
	"context"
	"fmt"
	repoOrder "order/internal/repository/order"
	platformMigrator "platform/pkg/migrator/pg"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func (s *OrderTestSuite) SetupSuite() {
	s.ctx = context.Background()

	pgWait := wait.ForLog("database system is ready to accept connections")
	pgWait.Occurrence = 2

	pgReq := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test_user",
			"POSTGRES_PASSWORD": "test_password",
			"POSTGRES_DB":       "test_orders_db",
		},
		WaitingFor: pgWait,
	}

	container, err := testcontainers.GenericContainer(s.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: pgReq,
		Started:          true,
	})
	s.Require().NoError(err)
	s.pgContainer = container

	pgHost, err := container.Host(s.ctx)
	s.Require().NoError(err)
	pgPort, err := container.MappedPort(s.ctx, "5432")
	s.Require().NoError(err)

	testDSN := fmt.Sprintf("postgres://test_user:test_password@%s:%s/test_orders_db?sslmode=disable", pgHost, pgPort.Port())

	dbMigrator := platformMigrator.NewMigrator(testDSN, "../../migrations")
	err = dbMigrator.Up()
	s.Require().NoError(err, "Не удалось накатать миграции в Testcontainers")

	pool, err := pgxpool.New(s.ctx, testDSN)
	s.Require().NoError(err)
	s.dbPool = pool

	s.repo = repoOrder.NewOrderRepository(s.dbPool)
}
