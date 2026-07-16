package env

import (
	"errors"
	"os"
)

type orderAssembledConsumerConfig struct {
	assembledTopic string
	groupID        string
}

func NewOrderAssembledConsumerConfig() (*orderAssembledConsumerConfig, error) {
	topic := os.Getenv("KAFKA_SHIP_ASSEMBLED_TOPIC")
	groupID := os.Getenv("KAFKA_SHIP_ASSEMBLED_GROUP_ID")
	if topic == "" || groupID == "" {
		return nil, errors.New("KAFKA_SHIP_ASSEMBLED_TOPIC or KAFKA_SHIP_ASSEMBLED_GROUP_ID is not set")
	}
	return &orderAssembledConsumerConfig{assembledTopic: topic, groupID: groupID}, nil
}

func (c *orderAssembledConsumerConfig) AssembledTopic() string { return c.assembledTopic }
func (c *orderAssembledConsumerConfig) GroupID() string        { return c.groupID }
