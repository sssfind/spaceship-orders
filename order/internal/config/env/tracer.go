package env

import "os"

type TracerConfig interface {
	CollectorEndpoint() string
	ServiceName() string
	Environment() string
	ServiceVersion() string
}

type tracerConfig struct {
	collectorEndpoint string
	serviceName       string
	environment       string
	serviceVersion    string
}

func NewTracerConfig() (TracerConfig, error) {
	collectorEndpoint := os.Getenv("OTEL_COLLECTOR_ENDPOINT")
	if collectorEndpoint == "" {
		collectorEndpoint = "otel-collector:4317"
	}

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "order-service"
	}

	environment := os.Getenv("APP_ENV")
	if environment == "" {
		environment = "local"
	}

	return &tracerConfig{
		collectorEndpoint: collectorEndpoint,
		serviceName:       serviceName,
		environment:       environment,
		serviceVersion:    "1.0.0",
	}, nil
}

func (c *tracerConfig) CollectorEndpoint() string { return c.collectorEndpoint }
func (c *tracerConfig) ServiceName() string       { return c.serviceName }
func (c *tracerConfig) Environment() string       { return c.environment }
func (c *tracerConfig) ServiceVersion() string    { return c.serviceVersion }
