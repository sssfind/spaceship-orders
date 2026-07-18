package app

import (
	"context"
	"fmt"
	"net"
	"platform/pkg/logger"
	platformMigrator "platform/pkg/migrator/pg"

	"iam/internal/config"
	authv1 "spaceship-orders/shared/pkg/proto/auth/v1"
	userv1 "spaceship-orders/shared/pkg/proto/user/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	config     *config.Config
	container  *Container
	grpcServer *grpc.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, fmt.Errorf("failed to init dependencies: %w", err)
	}

	return a, nil
}

func (a *App) Run() error {
	lis, err := net.Listen("tcp", a.config.GRPC().Address())
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", a.config.GRPC().Address(), err)
	}

	fmt.Printf("gRPC server is running on %s\n", a.config.GRPC().Address())

	if err := a.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC: %w", err)
	}

	return nil
}

func (a *App) Stop() error {
	fmt.Println("Stopping gRPC server gracefully...")
	a.grpcServer.GracefulStop()

	if err := a.container.Close(); err != nil {
		return fmt.Errorf("failed to close container: %w", err)
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	a.config = cfg

	err = logger.Init(cfg.Logger().Level(), cfg.Logger().AsJSON())
	if err != nil {
		return fmt.Errorf("failed to init logger: %w", err)
	}

	dbMigrator := platformMigrator.NewMigrator(
		cfg.Postgres().DSN(),
		cfg.Postgres().MigrationDir(),
	)

	if err := dbMigrator.Up(); err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}
	logger.Info(ctx, "Database migrations successfully applied")

	a.container = NewContainer(cfg)

	if err := a.initGRPCServer(ctx); err != nil {
		return err
	}

	return nil
}
func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer()

	reflection.Register(a.grpcServer)

	authImpl, err := a.container.AuthImpl(ctx)
	if err != nil {
		return err
	}

	userImpl, err := a.container.UserImpl(ctx)
	if err != nil {
		return err
	}

	authv1.RegisterAuthServiceServer(a.grpcServer, authImpl)
	userv1.RegisterUserServiceServer(a.grpcServer, userImpl)

	return nil
}
