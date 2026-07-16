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
	"platform/pkg/logger" // Твой пакет логгера[cite: 15]
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

	// 2. Регистрируем закрытие ресурсов через AddNamed.
	// Благодаря этому ты увидишь в консоли аккуратные отчеты по каждому ресурсу[cite: 18].
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

	// 3. Получаем и запускаем консьюмер в фоновой горутине
	cons, err := a.serviceProvider.OrderConsumer()
	if err != nil {
		return err
	}

	go func() {
		logger.Info(ctx, "Assembly Service Kafka Consumer is starting...")
		if err := cons.Run(ctx); err != nil {
			logger.Error(ctx, fmt.Sprintf("Assembly Consumer stopped with error: %v", err))
		}
	}()

	// 4. Блокируем выполнение и ждем сигналы ОС прямо здесь
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	logger.Info(ctx, fmt.Sprintf("Получен системный сигнал %v. Начинаем graceful shutdown...", sig))

	// 5. Вызываем закрытие ресурсов синхронно с таймаутом в 5 секунд
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := closer.CloseAll(shutdownCtx); err != nil {
		logger.Error(ctx, fmt.Sprintf("Ошибка при закрытии ресурсов приложения: %v", err))
		return err
	}

	logger.Info(ctx, "Assembly Service успешно остановлен.")
	return nil
}
