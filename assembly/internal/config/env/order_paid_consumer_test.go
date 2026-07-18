package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOrderPaidConsumerConfig_Success(t *testing.T) {
	t.Setenv("KAFKA_ORDER_PAID_TOPIC", "test-topic")
	t.Setenv("KAFKA_ORDER_PAID_GROUP_ID", "test-group")

	cfg, err := NewOrderPaidConsumerConfig()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "test-topic", cfg.PaidTopic())
	assert.Equal(t, "test-group", cfg.GroupID())
}

func TestNewOrderPaidConsumerConfig_MissingTopic(t *testing.T) {
	// Очищаем топик, но оставляем group ID
	t.Setenv("KAFKA_ORDER_PAID_TOPIC", "")
	t.Setenv("KAFKA_ORDER_PAID_GROUP_ID", "test-group")

	cfg, err := NewOrderPaidConsumerConfig()

	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.EqualError(t, err, "KAFKA_ORDER_PAID_TOPIC or KAFKA_ORDER_PAID_GROUP_ID is not set")
}

func TestNewOrderPaidConsumerConfig_MissingGroupID(t *testing.T) {
	t.Setenv("KAFKA_ORDER_PAID_TOPIC", "test-topic")
	t.Setenv("KAFKA_ORDER_PAID_GROUP_ID", "")

	cfg, err := NewOrderPaidConsumerConfig()

	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.EqualError(t, err, "KAFKA_ORDER_PAID_TOPIC or KAFKA_ORDER_PAID_GROUP_ID is not set")
}
