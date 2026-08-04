package env

import (
	"os"
	"strconv"
	"strings"
)

type loggerConfig struct {
	logLevel              string
	logAsJSON             bool
	serviceName           string
	outputs               []string
	otelCollectorEndpoint string
}

func NewLoggerConfig() (*loggerConfig, error) {
	level := os.Getenv("LOGGER_LEVEL")
	if level == "" {
		level = "info"
	}

	asJSONStr := os.Getenv("LOGGER_AS_JSON")
	asJSON, _ := strconv.ParseBool(asJSONStr)

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "IAM-service"
	}

	outputsStr := os.Getenv("LOG_OUTPUTS")
	var outputs []string
	if outputsStr != "" {
		for _, out := range strings.Split(outputsStr, ",") {
			if trimmed := strings.TrimSpace(out); trimmed != "" {
				outputs = append(outputs, trimmed)
			}
		}
	}
	if len(outputs) == 0 {
		outputs = []string{"stdout"}
	}

	otelEndpoint := os.Getenv("OTEL_COLLECTOR_ENDPOINT")
	if otelEndpoint == "" {
		otelEndpoint = "otel-collector:4317"
	}

	return &loggerConfig{
		logLevel:              level,
		logAsJSON:             asJSON,
		serviceName:           serviceName,
		outputs:               outputs,
		otelCollectorEndpoint: otelEndpoint,
	}, nil
}

func (l *loggerConfig) Level() string {
	return l.logLevel
}

func (l *loggerConfig) AsJSON() bool {
	return l.logAsJSON
}
func (cfg *loggerConfig) ServiceName() string           { return cfg.serviceName }
func (cfg *loggerConfig) Outputs() []string             { return cfg.outputs }
func (cfg *loggerConfig) OtelCollectorEndpoint() string { return cfg.otelCollectorEndpoint }
