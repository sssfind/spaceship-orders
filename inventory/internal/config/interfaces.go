package config

type LoggerConfig interface {
	GetLogLevel() string
	GetLogAsJSON() bool
}

type InventoryGrpcConfig interface {
	GetAddress() string
}

type MongoConfig interface {
	GetURI() string
	GetDatabaseName() string
}
