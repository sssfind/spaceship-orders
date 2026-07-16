package config

type LoggerConfig interface {
	LogLevel() string
	LogAsJSON() bool
}

type PaymentGrpcConfig interface {
	Address() string
}
