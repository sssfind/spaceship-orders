package config

import "time"

type LoggerConfig interface {
	LogLevel() string
	LogAsJSON() bool
	ServiceName() string
	Outputs() []string
	OtelCollectorEndpoint() string
}

type HTTPConfig interface {
	Address() string
	ReadTimeout() time.Duration
}

type MetricsHTTPConfig interface {
	Address() string
}

type DemoFaultConfig interface {
	Mode() string
}

type PostgresConfig interface {
	Dsn() string
	MigrationDir() string
}

type SessionConfig interface {
	TTL() time.Duration
}
