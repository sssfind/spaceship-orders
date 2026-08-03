package config

type LoggerConfig interface {
	LogLevel() string
	LogAsJSON() bool
	ServiceName() string
	Outputs() []string
	OtelCollectorEndpoint() string
}

type InventoryGrpcConfig interface {
	Address() string
}

type MongoConfig interface {
	GetURI() string
	GetDatabaseName() string
}
