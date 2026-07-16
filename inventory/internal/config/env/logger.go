package env

import (
	"os"
	"strconv"
)

type loggerConfig struct {
	level  string
	asJSON bool
}

func NewLoggerConfig() (*loggerConfig, error) {
	level := os.Getenv("LOGGER_LEVEL")
	if level == "" {
		level = "info"
	}

	asJSONStr := os.Getenv("LOGGER_AS_JSON")
	asJSON, _ := strconv.ParseBool(asJSONStr)

	return &loggerConfig{
		level:  level,
		asJSON: asJSON,
	}, nil
}

func (cfg *loggerConfig) LogLevel() string { return cfg.level }
func (cfg *loggerConfig) LogAsJSON() bool  { return cfg.asJSON }
