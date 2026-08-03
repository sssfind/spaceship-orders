package config

type LoggerConfig interface {
	LogLevel() string
	LogAsJSON() bool
	ServiceName() string
	Outputs() []string
	OtelCollectorEndpoint() string
}

type PaymentGrpcConfig interface {
	Address() string
}
