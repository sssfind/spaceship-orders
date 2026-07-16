package env

import (
	"errors"
	"os"
)

type orderAssembledProducerConfig struct {
	assembledTopic string
}

func NewOrderAssembledProducerConfig() (*orderAssembledProducerConfig, error) {
	topic := os.Getenv("KAFKA_SHIP_ASSEMBLED_TOPIC")
	if topic == "" {
		return nil, errors.New("KAFKA_SHIP_ASSEMBLED_TOPIC is not set")
	}
	return &orderAssembledProducerConfig{assembledTopic: topic}, nil
}

func (c *orderAssembledProducerConfig) AssembledTopic() string { return c.assembledTopic }
