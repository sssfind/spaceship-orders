package env

import (
	"os"
	"time"
)

type sessionConfig struct {
	ttl time.Duration
}

func NewSessionConfig() (*sessionConfig, error) {
	ttlStr := os.Getenv("SESSION_TTL")
	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		ttl = 24 * time.Hour
	}

	return &sessionConfig{ttl: ttl}, nil
}

func (s *sessionConfig) TTL() time.Duration {
	return s.ttl
}
