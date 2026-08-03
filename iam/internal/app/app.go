package app

import (
	"context"
	"fmt"
	"net"
	"syscall"

	"iam/internal/config"
	"platform/pkg/closer"
	"platform/pkg/logger"
	platformMigrator "platform/pkg/migrator/pg"
	authv1 "spaceship-orders/shared/pkg/proto/auth/v1"
	userv1 "spaceship-orders/shared/pkg/proto/user/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type App struct {
	cfg        *config.Config
	container  *Container
	grpcServer *grpc.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	a.cfg = cfg

	err = logger.InitWithConfig(logger.Config{
		Level:                 cfg.Logger().Level(),
		AsJSON:                cfg.Logger().AsJSON(),
		ServiceName:           cfg.Logger().ServiceName(),
		Outputs:               cfg.Logger().Outputs(),
		OtelCollectorEndpoint: cfg.Logger().OtelCollectorEndpoint(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init logger: %w", err)
	}

	closer.AddNamed("logger", func(c context.Context) error {
		return logger.Shutdown(c)
	})

	dbMigrator := platformMigrator.NewMigrator(
		cfg.Postgres().DSN(),
		cfg.Postgres().MigrationDir(),
	)
	if err := dbMigrator.Up(); err != nil {
		return nil, fmt.Errorf("database migration failed: %w", err)
	}
	logger.Info(ctx, "Database migrations successfully applied")

	a.container = NewContainer(cfg)
	closer.AddNamed("container", func(_ context.Context) error {
		return a.container.Close()
	})

	err = a.initGrpcServer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to init grpc server: %w", err)
	}

	return a, nil
}

func (a *App) initGrpcServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(a.grpcServer)

	// Настройка gRPC Health Check для двух сервисов IAM
	healthCheckServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(a.grpcServer, healthCheckServer)
	healthCheckServer.SetServingStatus("auth.v1.AuthService", grpc_health_v1.HealthCheckResponse_SERVING)
	healthCheckServer.SetServingStatus("user.v1.UserService", grpc_health_v1.HealthCheckResponse_SERVING)

	authImpl, err := a.container.AuthImpl(ctx)
	if err != nil {
		return fmt.Errorf("failed to init auth impl: %w", err)
	}

	userImpl, err := a.container.UserImpl(ctx)
	if err != nil {
		return fmt.Errorf("failed to init user impl: %w", err)
	}

	authv1.RegisterAuthServiceServer(a.grpcServer, authImpl)
	userv1.RegisterUserServiceServer(a.grpcServer, userImpl)

	closer.AddNamed("grpc_server", func(_ context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	return nil
}

func (a *App) Run() error {
	closer.Configure(syscall.SIGINT, syscall.SIGTERM)
	closer.SetLogger(logger.Logger())

	address := a.cfg.GRPC().Address()
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", address, err)
	}

	logger.Info(context.Background(), fmt.Sprintf("IAM gRPC Server успешно запущен на %s", address))

	if err := a.grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
		return fmt.Errorf("grpc server failure: %w", err)
	}

	return nil
}
