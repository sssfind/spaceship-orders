package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"assembly/internal/config"
	"platform/pkg/closer"
	"platform/pkg/logger"
)

type App struct {
	serviceProvider *serviceProvider
}

func NewApp(ctx context.Context) (*App, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to init config: %w", err)
	}

	err = logger.Init(cfg.LogLevel(), true)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize global logger: %w", err)
	}

	return &App{
		serviceProvider: newServiceProvider(cfg),
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	closer.SetLogger(logger.Logger())

	closer.AddNamed("Logger Sync", func(_ context.Context) error {
		return logger.Sync()
	})

	closer.AddNamed("Sarama Producer", func(_ context.Context) error {
		if a.serviceProvider.saramaProducer != nil {
			return a.serviceProvider.saramaProducer.Close()
		}
		return nil
	})

	closer.AddNamed("Sarama Consumer Group", func(_ context.Context) error {
		if a.serviceProvider.saramaConsumer != nil {
			return a.serviceProvider.saramaConsumer.Close()
		}
		return nil
	})

	cons, err := a.serviceProvider.OrderConsumer()
	if err != nil {
		return err
	}
	cfg := a.serviceProvider.cfg
	fmt.Fprintf(os.Stderr, "!!! DEBUG: Brokers=%v, GroupID=%s, Topic=%s\n", cfg.Brokers(), cfg.GroupID(), cfg.PaidTopic())

	go func() {
		logger.Info(ctx, "Assembly Service Kafka Consumer is starting...")
		if err := cons.Run(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "!!! CRITICAL CONSUMER ERROR: %v\n", err)
			logger.Error(ctx, fmt.Sprintf("Assembly Consumer stopped with error: %v", err))
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	logger.Info(ctx, fmt.Sprintf("Получен системный сигнал %v. Начинаем graceful shutdown...", sig))

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := closer.CloseAll(shutdownCtx); err != nil {
		logger.Error(ctx, fmt.Sprintf("Ошибка при закрытии ресурсов приложения: %v", err))
		return err
	}

	logger.Info(ctx, "Assembly Service успешно остановлен.")
	return nil
}
