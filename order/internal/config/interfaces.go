package config

import "time"

type LoggerConfig interface {
	LogLevel() string
	LogAsJSON() bool
	ServiceName() string
	Outputs() []string
	OtelCollectorEndpoint() string
}

type TracerConfig interface {
	CollectorEndpoint() string
	ServiceName() string
	Environment() string
	ServiceVersion() string
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
