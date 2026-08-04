package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"syscall"

	"assembly/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"platform/pkg/closer"
	"platform/pkg/logger"
)

type App struct {
	serviceProvider *serviceProvider
	httpServer      *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	cfg, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to init config: %w", err)
	}

	err = logger.InitWithConfig(logger.Config{
		Level:                 cfg.LogLevel(),
		AsJSON:                cfg.LogAsJSON(),
		ServiceName:           cfg.ServiceName(),
		Outputs:               cfg.Outputs(),
		OtelCollectorEndpoint: cfg.OtelCollectorEndpoint(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize global logger: %w", err)
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
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Эндпоинт для Prometheus
	r.Handle("/metrics", promhttp.Handler())

	a.httpServer = &http.Server{
		Addr:    a.serviceProvider.cfg.Address(), // Берем HTTP-адрес из конфига (например, :8081)
		Handler: r,
	}

	closer.AddNamed("http_server", func(c context.Context) error {
		return a.httpServer.Shutdown(c)
	})

	return nil
}

func (a *App) Run(ctx context.Context) error {
	closer.Configure(syscall.SIGINT, syscall.SIGTERM)
	closer.SetLogger(logger.Logger())

	cons, err := a.serviceProvider.OrderConsumer()
	if err != nil {
		return fmt.Errorf("failed to init order consumer: %w", err)
	}

	// 1. Запускаем Kafka Consumer в фоновой горутине
	go func() {
		logger.Info(ctx, fmt.Sprintf("Assembly Service Kafka Consumer успешно запущен на брокерах %v", a.serviceProvider.cfg.Brokers()))
		if err := cons.Run(ctx); err != nil {
			logger.Error(ctx, fmt.Sprintf("Assembly Consumer stopped with error: %v", err))
		}
	}()

	// 2. HTTP Server запускаем на основном потоке
	logger.Info(ctx, fmt.Sprintf("Assembly HTTP Server успешно запущен на %s", a.httpServer.Addr))

	if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server runtime error: %w", err)
	}

	return nil
}
