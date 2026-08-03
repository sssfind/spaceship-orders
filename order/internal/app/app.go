package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"syscall"

	healthAPI "order/internal/api/health"
	"order/internal/config"
	customMiddleware "order/internal/middleware"
	"platform/pkg/closer"
	"platform/pkg/logger"
	platformMigrator "platform/pkg/migrator/pg"
	orderV1 "spaceship-orders/shared/pkg/openapi/order/v1"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	serviceProvider *serviceProvider
	httpServer      *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	// Загружаем конфигурацию из переменных окружения
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Инициализируем платформенный логгер с поддержкой мульти-кора
	err = logger.InitWithConfig(logger.Config{
		Level:                 cfg.LogLevel(),
		AsJSON:                cfg.LogAsJSON(),
		ServiceName:           cfg.ServiceName(),
		Outputs:               cfg.Outputs(),
		OtelCollectorEndpoint: cfg.OtelCollectorEndpoint(),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to init logger: %w", err)
	}

	closer.AddNamed("logger", func(c context.Context) error {
		return logger.Shutdown(c)
	})

	a.serviceProvider = newServiceProvider(cfg)

	err = a.initDependencies(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initDependencies(ctx context.Context) error {
	dbMigrator := platformMigrator.NewMigrator(
		a.serviceProvider.cfg.Dsn(),
		a.serviceProvider.cfg.MigrationDir(),
	)
	if err := dbMigrator.Up(); err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}
	logger.Info(ctx, "Database migrations successfully applied")

	// Собираем хендлеры
	apiHandler, err := a.serviceProvider.APIHandler(ctx)
	if err != nil {
		return err
	}

	orderServer, err := orderV1.NewServer(apiHandler)
	if err != nil {
		return fmt.Errorf("failed to create openapi server: %w", err)
	}

	healthHandler := healthAPI.NewHandler()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(customMiddleware.RequestLogger)

	r.Get("/health", healthHandler.HealthCheck)

	r.Mount("/", orderServer)

	a.httpServer = &http.Server{
		Addr:              a.serviceProvider.cfg.OrderHttpConfig.Address(),
		Handler:           r,
		ReadHeaderTimeout: a.serviceProvider.cfg.ReadTimeout(),
	}

	closer.AddNamed("http_server", func(c context.Context) error {
		return a.httpServer.Shutdown(c)
	})

	return nil
}

func (a *App) Run() error {
	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	closer.SetLogger(logger.Logger())

	ctx := context.Background()

	cons, err := a.serviceProvider.OrderConsumer(ctx)
	if err != nil {
		return fmt.Errorf("failed to run order consumer: %w", err)
	}

	go func() {
		logger.Info(ctx, "Order Service Kafka Consumer is starting...")
		if err := cons.Run(ctx); err != nil {
			logger.Error(ctx, fmt.Sprintf("Order Consumer stopped with error: %v", err))
		}
	}()

	logger.Info(ctx, fmt.Sprintf("order HTTP Server успешно запущен на %s", a.httpServer.Addr))

	if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server runtime error: %w", err)
	}

	return nil
}
