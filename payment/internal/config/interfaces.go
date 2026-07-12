package config

type LoggerConfig interface {
	GetLogLevel() string
	GetLogAsJSON() bool
}

type PaymentGrpcConfig interface {
	GetAddress() string
}
