package config

import (
	"fmt"

	"assembly/internal/config/env"
)

type cfg struct {
	kafkaConfig                  KafkaConfig
	loggerConfig                 LoggerConfig
	orderPaidConsumerConfig      OrderPaidConsumerConfig
	orderAssembledProducerConfig OrderAssembledProducerConfig
	httpConfig                   HTTPConfig
}

func NewConfig() (Config, error) {
	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load kafka config: %w", err)
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load logger config: %w", err)
	}

	paidConsumerCfg, err := env.NewOrderPaidConsumerConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load order paid consumer config: %w", err)
	}

	assembledProducerCfg, err := env.NewOrderAssembledProducerConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load order assembled producer config: %w", err)
	}

	httpCfg, err := env.NewHTTPConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load http config: %w", err)
	}

	return &cfg{
		kafkaConfig:                  kafkaCfg,
		loggerConfig:                 loggerCfg,
		orderPaidConsumerConfig:      paidConsumerCfg,
		orderAssembledProducerConfig: assembledProducerCfg,
		httpConfig:                   httpCfg,
	}, nil
}

func (c *cfg) Brokers() []string             { return c.kafkaConfig.Brokers() }
func (c *cfg) PaidTopic() string             { return c.orderPaidConsumerConfig.PaidTopic() }
func (c *cfg) GroupID() string               { return c.orderPaidConsumerConfig.GroupID() }
func (c *cfg) AssembledTopic() string        { return c.orderAssembledProducerConfig.AssembledTopic() }
func (c *cfg) LogLevel() string              { return c.loggerConfig.LogLevel() }
func (c *cfg) LogAsJSON() bool               { return c.loggerConfig.LogAsJSON() }
func (c *cfg) ServiceName() string           { return c.loggerConfig.ServiceName() }
func (c *cfg) Outputs() []string             { return c.loggerConfig.Outputs() }
func (c *cfg) OtelCollectorEndpoint() string { return c.loggerConfig.OtelCollectorEndpoint() }
func (c *cfg) Address() string               { return c.httpConfig.Address() }
