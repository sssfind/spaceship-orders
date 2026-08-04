package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"payment/internal/config/env"
)

type Config struct {
	LoggerConfig
	TracerConfig
	PaymentGrpcConfig
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return nil, fmt.Errorf("logger config error: %w", err)
	}

	tracerCfg, err := env.NewTracerConfig()
	if err != nil {
		return nil, fmt.Errorf("tracer config error: %w", err)
	}

	paymentCfg, err := env.NewPaymentGrpcConfig()
	if err != nil {
		return nil, fmt.Errorf("payment config error: %w", err)
	}

	return &Config{
		LoggerConfig:      loggerCfg,
		TracerConfig:      tracerCfg,
		PaymentGrpcConfig: paymentCfg,
	}, nil
}

func (c *Config) ServiceName() string {
	return c.LoggerConfig.ServiceName()
}
