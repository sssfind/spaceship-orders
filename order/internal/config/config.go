package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"order/internal/config/env"
)

type Config struct {
	LoggerConfig
	InventoryGrpcConfig
	PaymentGrpcConfig
	OrderHttpConfig
	PostgresConfig
	KafkaConfig
	OrderPaidProducerConfig
	OrderAssembledConsumerConfig
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

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return nil, fmt.Errorf("kafka config error: %w", err)
	}

	paidProducerCfg, err := env.NewOrderPaidProducerConfig()
	if err != nil {
		return nil, fmt.Errorf("order paid producer config error: %w", err)
	}

	assembledConsumerCfg, err := env.NewOrderAssembledConsumerConfig()
	if err != nil {
		return nil, fmt.Errorf("order assembled consumer config error: %w", err)
	}

	return &Config{
		LoggerConfig:                 loggerCfg,
		InventoryGrpcConfig:          inventoryCfg,
		PaymentGrpcConfig:            paymentCfg,
		OrderHttpConfig:              httpCfg,
		PostgresConfig:               postgresCfg,
		KafkaConfig:                  kafkaCfg,
		OrderPaidProducerConfig:      paidProducerCfg,
		OrderAssembledConsumerConfig: assembledConsumerCfg,
	}, nil
}
