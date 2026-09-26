package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	apiV1 "order/internal/api/order/v1"
	"order/internal/config"
	"order/internal/repository"
	repoOrder "order/internal/repository/order"
	repoPart "order/internal/repository/part"
	repoPayment "order/internal/repository/payment"
	"order/internal/service"
	orderImpl "order/internal/service/order"
	partImpl "order/internal/service/part"
	"platform/pkg/closer"

	orderV1 "spaceship-orders/shared/pkg/openapi/order/v1"
)

type serviceProvider struct {
	cfg *config.Config

	dbPool      *pgxpool.Pool
	orderRepo   repository.OrderRepository
	partRepo    repository.PartRepository
	paymentRepo repository.PaymentRepository
	orderSrv    service.OrderService
	partSrv     service.PartService
	apiHandler  orderV1.Handler
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

func (sp *serviceProvider) PartRepository(ctx context.Context) (repository.PartRepository, error) {
	if sp.partRepo == nil {
		pool, err := sp.DBPool(ctx)
		if err != nil {
			return nil, err
		}
		sp.partRepo = repoPart.NewPartRepository(pool)
	}
	return sp.partRepo, nil
}

func (sp *serviceProvider) PaymentRepository(ctx context.Context) (repository.PaymentRepository, error) {
	if sp.paymentRepo == nil {
		pool, err := sp.DBPool(ctx)
		if err != nil {
			return nil, err
		}
		sp.paymentRepo = repoPayment.NewPaymentRepository(pool)
	}
	return sp.paymentRepo, nil
}

func (sp *serviceProvider) OrderService(ctx context.Context) (service.OrderService, error) {
	if sp.orderSrv == nil {
		orderRepo, err := sp.OrderRepository(ctx)
		if err != nil {
			return nil, err
		}
		partRepo, err := sp.PartRepository(ctx)
		if err != nil {
			return nil, err
		}
		paymentRepo, err := sp.PaymentRepository(ctx)
		if err != nil {
			return nil, err
		}
		sp.orderSrv = orderImpl.NewService(orderRepo, partRepo, paymentRepo)
	}
	return sp.orderSrv, nil
}

func (sp *serviceProvider) PartService(ctx context.Context) (service.PartService, error) {
	if sp.partSrv == nil {
		partRepo, err := sp.PartRepository(ctx)
		if err != nil {
			return nil, err
		}
		sp.partSrv = partImpl.NewService(partRepo)
	}
	return sp.partSrv, nil
}

func (sp *serviceProvider) APIHandler(ctx context.Context) (orderV1.Handler, error) {
	if sp.apiHandler == nil {
		orderSrv, err := sp.OrderService(ctx)
		if err != nil {
			return nil, err
		}
		partSrv, err := sp.PartService(ctx)
		if err != nil {
			return nil, err
		}
		sp.apiHandler = apiV1.NewAPI(orderSrv, partSrv)
	}
	return sp.apiHandler, nil
}
