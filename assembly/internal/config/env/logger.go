package env

import "os"

type loggerConfig struct {
	logLevel string
}

func NewLoggerConfig() (*loggerConfig, error) {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}
	return &loggerConfig{logLevel: level}, nil
}

func (c *loggerConfig) LogLevel() string { return c.logLevel }
