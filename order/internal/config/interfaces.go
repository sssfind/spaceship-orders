package config

import "time"

type LoggerConfig interface {
	GetLogLevel() string
	GetLogAsJSON() bool
}

type OrderHttpConfig interface {
	GetAddress() string
	GetReadTimeout() time.Duration
}

type PostgresConfig interface {
	GetDSN() string
	GetMigrationDir() string
}

type InventoryGrpcConfig interface {
	GetAddress() string
}

type PaymentGrpcConfig interface {
	GetAddress() string
}
