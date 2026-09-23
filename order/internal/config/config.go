package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"order/internal/config/env"
)

type Config struct {
	LoggerConfig
	TracerConfig
	OrderHttpConfig
	PostgresConfig
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

	httpCfg, err := env.NewOrderHTTPConfig()
	if err != nil {
		return nil, fmt.Errorf("http config error: %w", err)
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return nil, fmt.Errorf("postgres config error: %w", err)
	}

	return &Config{
		LoggerConfig:    loggerCfg,
		TracerConfig:    tracerCfg,
		OrderHttpConfig: httpCfg,
		PostgresConfig:  postgresCfg,
	}, nil
}

func (c *Config) ServiceName() string {
	return c.LoggerConfig.ServiceName()
}
