package config

type KafkaConfig interface {
	Brokers() []string
}

type LoggerConfig interface {
	LogLevel() string
}

type OrderPaidConsumerConfig interface {
	PaidTopic() string
	GroupID() string
}

type OrderAssembledProducerConfig interface {
	AssembledTopic() string
}

type Config interface {
	KafkaConfig
	LoggerConfig
	OrderPaidConsumerConfig
	OrderAssembledProducerConfig
}
