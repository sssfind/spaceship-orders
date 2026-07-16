package env

import (
	"errors"
	"os"
)

type orderPaidProducerConfig struct {
	paidTopic string
}

func NewOrderPaidProducerConfig() (*orderPaidProducerConfig, error) {
	topic := os.Getenv("KAFKA_ORDER_PAID_TOPIC")
	if topic == "" {
		return nil, errors.New("KAFKA_ORDER_PAID_TOPIC is not set")
	}
	return &orderPaidProducerConfig{paidTopic: topic}, nil
}

func (c *orderPaidProducerConfig) PaidTopic() string { return c.paidTopic }
