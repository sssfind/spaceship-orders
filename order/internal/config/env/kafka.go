package env

import (
	"errors"
	"os"
	"strings"
)

type kafkaConfig struct {
	brokers []string
}

func NewKafkaConfig() (*kafkaConfig, error) {
	brokersStr := os.Getenv("KAFKA_BROKERS")
	if brokersStr == "" {
		return nil, errors.New("KAFKA_BROKERS is not set")
	}
	return &kafkaConfig{brokers: strings.Split(brokersStr, ",")}, nil
}

func (c *kafkaConfig) Brokers() []string { return c.brokers }
