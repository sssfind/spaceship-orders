package app

import (
	"context"
	"fmt"
	"inventory/internal/config"
	"net"
	"syscall"

	apiV1 "inventory/internal/api/inventory/v1"
	"platform/pkg/closer"
	"platform/pkg/logger"
	pbInventory "spaceship-orders/shared/pkg/proto/inventory/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	err = logger.Init(cfg.GetLogLevel(), cfg.GetLogAsJSON())
	if err != nil {
		return nil, fmt.Errorf("failed to init logger: %w", err)
	}

	a.serviceProvider = newServiceProvider(cfg)

	err = a.initGrpcServer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to init grpc server: %w", err)
	}

	return a, nil
}

func (a *App) initGrpcServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	healthCheckServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(a.grpcServer, healthCheckServer)
	healthCheckServer.SetServingStatus("inventory.v1.InventoryService", grpc_health_v1.HealthCheckResponse_SERVING)

	srv, err := a.serviceProvider.PartService(ctx)
	if err != nil {
		return err
	}
	apiHandler := apiV1.NewAPI(srv)
	pbInventory.RegisterInventoryServiceServer(a.grpcServer, apiHandler)

	closer.AddNamed("grpc_server", func(_ context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	return nil
}

func (a *App) Run() error {
	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	lis, err := net.Listen("tcp", a.serviceProvider.cfg.GetAddress())
	if err != nil {
		return fmt.Errorf("failed to listen port: %w", err)
	}

	logger.Info(context.Background(), fmt.Sprintf("Inventory gRPC Server успешно запущен на %s", a.serviceProvider.cfg.GetAddress()))

	if err := a.grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
		return fmt.Errorf("grpc server failure: %w", err)
	}

	return nil
}
