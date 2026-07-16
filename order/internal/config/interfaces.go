package config

import "time"

type LoggerConfig interface {
	LogLevel() string
	LogAsJSON() bool
}

type OrderHttpConfig interface {
	Address() string
	ReadTimeout() time.Duration
}

type PostgresConfig interface {
	Dsn() string
	MigrationDir() string
}

type InventoryGrpcConfig interface {
	Address() string
}

type PaymentGrpcConfig interface {
	Address() string
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidProducerConfig interface {
	PaidTopic() string
}

type OrderAssembledConsumerConfig interface {
	AssembledTopic() string
	GroupID() string
}
