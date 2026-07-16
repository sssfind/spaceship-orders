package app

import (
	"context"
	"fmt"
	"net"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection" // Импортируем пакет рефлексии
	"payment/internal/config"
	"platform/pkg/closer"
	"platform/pkg/logger"
	pb "spaceship-orders/shared/pkg/proto/payment/v1"
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

	err = logger.Init(cfg.LogLevel(), cfg.LogAsJSON())
	if err != nil {
		return nil, fmt.Errorf("failed to init logger: %w", err)
	}

	a.serviceProvider = newServiceProvider(cfg)

	err = a.initGrpcServer(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initGrpcServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(a.grpcServer, healthServer)
	healthServer.SetServingStatus("payment.v1.PaymentService", grpc_health_v1.HealthCheckResponse_SERVING)

	apiHandler, err := a.serviceProvider.APIHandler()
	if err != nil {
		return err
	}
	pb.RegisterPaymentServiceServer(a.grpcServer, apiHandler)

	reflection.Register(a.grpcServer)

	// Graceful shutdown
	closer.AddNamed("payment_grpc_server", func(_ context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	return nil
}

func (a *App) Run() error {
	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	lis, err := net.Listen("tcp", a.serviceProvider.cfg.Address())
	if err != nil {
		return fmt.Errorf("failed to listen port: %w", err)
	}

	logger.Info(context.Background(), fmt.Sprintf("Payment gRPC Server успешно запущен на %s", a.serviceProvider.cfg.Address()))

	if err := a.grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
		return fmt.Errorf("grpc server failure: %w", err)
	}

	return nil
}
