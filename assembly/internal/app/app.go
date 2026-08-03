package app

import (
	"context"
	"fmt"
	"syscall"

	"assembly/internal/config"
	"platform/pkg/closer"
	"platform/pkg/logger"
)

type App struct {
	serviceProvider *serviceProvider
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

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	closer.Configure(syscall.SIGINT, syscall.SIGTERM)
	closer.SetLogger(logger.Logger())

	cons, err := a.serviceProvider.OrderConsumer()
	if err != nil {
		return fmt.Errorf("failed to init order consumer: %w", err)
	}

	logger.Info(ctx, fmt.Sprintf("Assembly Service Kafka Consumer успешно запущен на брокерах %v", a.serviceProvider.cfg.Brokers()))

	if err := cons.Run(ctx); err != nil {
		return fmt.Errorf("assembly consumer runtime error: %w", err)
	}

	return nil
}
