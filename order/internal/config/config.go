package config

import (
	"fmt"
	"order/internal/config/env"

	"github.com/joho/godotenv"
)

type Config struct {
	LoggerConfig
	InventoryGrpcConfig
	PaymentGrpcConfig
	OrderHttpConfig
	PostgresConfig
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return nil, fmt.Errorf("logger config error: %w", err)
	}

	inventoryCfg, err := env.NewInventoryGrpcConfig()
	if err != nil {
		return nil, fmt.Errorf("inventory config error: %w", err)
	}

	paymentCfg, err := env.NewPaymentGrpcConfig()
	if err != nil {
		return nil, fmt.Errorf("payment config error: %w", err)
	}

	httpCfg, err := env.NewOrderHTTPConfig()
	if err != nil {
		return nil, fmt.Errorf("http config error: %w", err)
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return nil, fmt.Errorf("postgres config error: %w", err)
	}

	return &Config{
		LoggerConfig:        loggerCfg,
		InventoryGrpcConfig: inventoryCfg,
		PaymentGrpcConfig:   paymentCfg,
		OrderHttpConfig:     httpCfg,
		PostgresConfig:      postgresCfg,
	}, nil
}
