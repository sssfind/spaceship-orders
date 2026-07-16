package env

import (
	"errors"
	"os"
)

type orderPaidConsumerConfig struct {
	paidTopic string
	groupID   string
}

func NewOrderPaidConsumerConfig() (*orderPaidConsumerConfig, error) {
	topic := os.Getenv("KAFKA_ORDER_PAID_TOPIC")
	groupID := os.Getenv("KAFKA_ORDER_PAID_GROUP_ID")
	if topic == "" || groupID == "" {
		return nil, errors.New("KAFKA_ORDER_PAID_TOPIC or KAFKA_ORDER_PAID_GROUP_ID is not set")
	}
	return &orderPaidConsumerConfig{paidTopic: topic, groupID: groupID}, nil
}

func (c *orderPaidConsumerConfig) PaidTopic() string { return c.paidTopic }
func (c *orderPaidConsumerConfig) GroupID() string   { return c.groupID }
