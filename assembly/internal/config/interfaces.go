package config

type KafkaConfig interface {
	Brokers() []string
}

type LoggerConfig interface {
	LogLevel() string
	LogAsJSON() bool
	ServiceName() string
	Outputs() []string
	OtelCollectorEndpoint() string
}

type OrderPaidConsumerConfig interface {
	PaidTopic() string
	GroupID() string
}

type OrderAssembledProducerConfig interface {
	AssembledTopic() string
}

type HTTPConfig interface {
	Address() string
}

type Config interface {
	KafkaConfig
	LoggerConfig
	OrderPaidConsumerConfig
	OrderAssembledProducerConfig
	HTTPConfig
}
