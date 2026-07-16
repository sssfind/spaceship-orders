package config

type LoggerConfig interface {
	LogLevel() string
	LogAsJSON() bool
}

type InventoryGrpcConfig interface {
	Address() string
}

type MongoConfig interface {
	GetURI() string
	GetDatabaseName() string
}
