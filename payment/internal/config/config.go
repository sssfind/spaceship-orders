package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"payment/internal/config/env"
)

type Config struct {
	LoggerConfig
	PaymentGrpcConfig
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return nil, fmt.Errorf("logger config error: %w", err)
	}

	paymentCfg, err := env.NewPaymentGrpcConfig()
	if err != nil {
		return nil, fmt.Errorf("mongo config error: %w", err)
	}

	return &Config{
		LoggerConfig:      loggerCfg,
		PaymentGrpcConfig: paymentCfg,
	}, nil
}
