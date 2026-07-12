package app

import (
	"context"
	"fmt"
	apiV1 "order/internal/api/order/v1"
	clientInventory "order/internal/client/grpc/inventory/v1"
	clientPayment "order/internal/client/grpc/payment/v1"
	"order/internal/config"
	"order/internal/repository"
	repoOrder "order/internal/repository/order"
	"order/internal/service"
	orderImpl "order/internal/service/order"
	"platform/pkg/closer"
	orderV1 "spaceship-orders/shared/pkg/openapi/order/v1"
	pbInventory "spaceship-orders/shared/pkg/proto/inventory/v1"
	pbPayment "spaceship-orders/shared/pkg/proto/payment/v1"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type serviceProvider struct {
	cfg *config.Config

	dbPool          *pgxpool.Pool
	invConn         *grpc.ClientConn
	payConn         *grpc.ClientConn
	inventoryClient orderImpl.InventoryClient
	paymentClient   orderImpl.PaymentClient
	orderRepo       repository.OrderRepository
	orderSrv        service.OrderService
	apiHandler      orderV1.Handler
}

func newServiceProvider(cfg *config.Config) *serviceProvider {
	return &serviceProvider{cfg: cfg}
}

// DBPool инициализирует пул соединений с PostgreSQL
func (sp *serviceProvider) DBPool(ctx context.Context) (*pgxpool.Pool, error) {
	if sp.dbPool == nil {
		pool, err := pgxpool.New(ctx, sp.cfg.GetDSN())
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

// InventoryGrpcClient настраивает gRPC-клиент к InventoryService
func (sp *serviceProvider) InventoryGrpcClient(ctx context.Context) (orderImpl.InventoryClient, error) {
	if sp.inventoryClient == nil {
		conn, err := grpc.NewClient(sp.cfg.InventoryGrpcConfig.GetAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("inventory connection error: %w", err)
		}
		sp.invConn = conn

		closer.AddNamed("grpc_inventory_connection", func(_ context.Context) error {
			return sp.invConn.Close()
		})

		pbClient := pbInventory.NewInventoryServiceClient(sp.invConn)
		sp.inventoryClient = clientInventory.NewClient(pbClient)
	}
	return sp.inventoryClient, nil
}

// PaymentGrpcClient настраивает gRPC-клиент к PaymentService
func (sp *serviceProvider) PaymentGrpcClient(ctx context.Context) (orderImpl.PaymentClient, error) {
	if sp.paymentClient == nil {
		conn, err := grpc.NewClient(sp.cfg.PaymentGrpcConfig.GetAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("payment connection error: %w", err)
		}
		sp.payConn = conn

		closer.AddNamed("grpc_payment_connection", func(_ context.Context) error {
			return sp.payConn.Close()
		})

		pbClient := pbPayment.NewPaymentServiceClient(sp.payConn)
		sp.paymentClient = clientPayment.NewClient(pbClient)
	}
	return sp.paymentClient, nil
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
		invClient, err := sp.InventoryGrpcClient(ctx)
		if err != nil {
			return nil, err
		}
		payClient, err := sp.PaymentGrpcClient(ctx)
		if err != nil {
			return nil, err
		}

		sp.orderSrv = orderImpl.NewService(repo, invClient, payClient)
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
