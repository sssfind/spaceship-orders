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
	asJSON, err := strconv.ParseBool(asJSONStr)
	if err != nil {
		asJSON = false
	}

	return &loggerConfig{level: level, asJSON: asJSON}, nil
}

func (l *loggerConfig) Level() string {
	return l.level
}

func (l *loggerConfig) AsJSON() bool {
	return l.asJSON
}
