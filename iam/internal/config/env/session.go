package env

import (
	"fmt"
	"os"
	"time"
)

type sessionConfig struct {
	ttl time.Duration
}

func NewSessionConfig() (*sessionConfig, error) {
	ttlStr := os.Getenv("SESSION_TTL")
	if ttlStr == "" {
		ttlStr = "24h"
	}

	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SESSION_TTL format: %w", err)
	}

	return &sessionConfig{ttl: ttl}, nil
}

func (cfg *sessionConfig) TTL() time.Duration { return cfg.ttl }
