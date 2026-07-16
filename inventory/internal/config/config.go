package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"inventory/internal/config/env"
)

type Config struct {
	LoggerConfig
	MongoConfig
	InventoryGrpcConfig
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return nil, fmt.Errorf("logger config error: %w", err)
	}

	mongoCfg, err := env.NewMongoConfig()
	if err != nil {
		return nil, fmt.Errorf("mongo config error: %w", err)
	}

	grpcCfg, err := env.NewInventoryGrpcConfig()
	if err != nil {
		return nil, fmt.Errorf("grpc config error: %w", err)
	}

	return &Config{
		LoggerConfig:        loggerCfg,
		MongoConfig:         mongoCfg,
		InventoryGrpcConfig: grpcCfg,
	}, nil
}
