package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	healthAPI "order/internal/api/health"
	"order/internal/config"
	customMiddleware "order/internal/middleware"
	"platform/pkg/closer"
	"platform/pkg/logger"
	platformMigrator "platform/pkg/migrator/pg"
	"syscall"

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

	// Инициализируем глобальный структурированный логгер платформы
	err = logger.Init(cfg.GetLogLevel(), cfg.GetLogAsJSON())
	if err != nil {
		return nil, fmt.Errorf("failed to init logger: %w", err)
	}

	a.serviceProvider = newServiceProvider(cfg)

	// Инициализируем базовую инфраструктуру и сервер
	err = a.initDependencies(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initDependencies(ctx context.Context) error {
	// Накатываем миграции БД через платформенную библиотеку
	dbMigrator := platformMigrator.NewMigrator(
		a.serviceProvider.cfg.GetDSN(),
		a.serviceProvider.cfg.GetMigrationDir(),
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
		Addr:              a.serviceProvider.cfg.OrderHttpConfig.GetAddress(),
		Handler:           r,
		ReadHeaderTimeout: a.serviceProvider.cfg.GetReadTimeout(),
	}

	closer.AddNamed("http_server", func(c context.Context) error {
		return a.httpServer.Shutdown(c)
	})

	return nil
}

func (a *App) Run() error {
	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	logger.Info(context.Background(), fmt.Sprintf("order HTTP Server успешно запущен на %s", a.httpServer.Addr))

	if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server runtime error: %w", err)
	}

	return nil
}
