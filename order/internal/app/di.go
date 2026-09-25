package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	apiV1 "order/internal/api/order/v1"
	"order/internal/config"
	"order/internal/repository"
	repoOrder "order/internal/repository/order"
	"order/internal/service"
	orderImpl "order/internal/service/order"
	"platform/pkg/closer"

	orderV1 "spaceship-orders/shared/pkg/openapi/order/v1"
)

type serviceProvider struct {
	cfg *config.Config

	dbPool     *pgxpool.Pool
	orderRepo  repository.OrderRepository
	orderSrv   service.OrderService
	apiHandler orderV1.Handler
}

func newServiceProvider(cfg *config.Config) *serviceProvider {
	return &serviceProvider{cfg: cfg}
}

func (sp *serviceProvider) DBPool(ctx context.Context) (*pgxpool.Pool, error) {
	if sp.dbPool == nil {
		pool, err := pgxpool.New(ctx, sp.cfg.Dsn())
		if err != nil {
			return nil, fmt.Errorf("failed to connect to postgres: %w", err)
		}
		sp.dbPool = pool

		closer.AddNamed("postgres_pool", func(_ context.Context) error {
			sp.dbPool.Close()
			return nil
		})
	}
	return sp.dbPool, nil
}

func (sp *serviceProvider) OrderRepository(ctx context.Context) (repository.OrderRepository, error) {
	if sp.orderRepo == nil {
		pool, err := sp.DBPool(ctx)
		if err != nil {
			return nil, err
		}
		sp.orderRepo = repoOrder.NewOrderRepository(pool)
	}
	return sp.orderRepo, nil
}

func (sp *serviceProvider) OrderService(ctx context.Context) (service.OrderService, error) {
	if sp.orderSrv == nil {
		repo, err := sp.OrderRepository(ctx)
		if err != nil {
			return nil, err
		}

		sp.orderSrv = orderImpl.NewService(repo)
	}
	return sp.orderSrv, nil
}

func (sp *serviceProvider) APIHandler(ctx context.Context) (orderV1.Handler, error) {
	if sp.apiHandler == nil {
		srv, err := sp.OrderService(ctx)
		if err != nil {
			return nil, err
		}
		sp.apiHandler = apiV1.NewAPI(srv)
	}
	return sp.apiHandler, nil
}
