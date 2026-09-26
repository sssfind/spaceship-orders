package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"iam/internal/config/env"
)

type Config struct {
	LoggerConfig
	HTTPConfig
	PostgresConfig
	SessionConfig
}

func Load() (*Config, error) {
	loadDotEnv()

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return nil, fmt.Errorf("logger config error: %w", err)
	}

	httpCfg, err := env.NewHTTPConfig()
	if err != nil {
		return nil, fmt.Errorf("http config error: %w", err)
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return nil, fmt.Errorf("postgres config error: %w", err)
	}

	sessionCfg, err := env.NewSessionConfig()
	if err != nil {
		return nil, fmt.Errorf("session config error: %w", err)
	}

	return &Config{
		LoggerConfig:   loggerCfg,
		HTTPConfig:     httpCfg,
		PostgresConfig: postgresCfg,
		SessionConfig:  sessionCfg,
	}, nil
}

func (c *Config) ServiceName() string {
	return c.LoggerConfig.ServiceName()
}

func loadDotEnv() {
	candidates := []string{".env"}

	if composeEnv, ok := findFileUpwards("deploy/compose/iam/.env"); ok {
		candidates = append(candidates, composeEnv)
	} else {
		candidates = append(candidates,
			filepath.Join("..", "deploy", "compose", "iam", ".env"),
		)
	}

	seen := make(map[string]struct{}, len(candidates))
	for _, path := range candidates {
		abs, err := filepath.Abs(path)
		if err != nil {
			abs = path
		}
		if _, ok := seen[abs]; ok {
			continue
		}
		seen[abs] = struct{}{}

		if _, err := os.Stat(path); err != nil {
			continue
		}
		_ = godotenv.Load(path)
	}
}

func findFileUpwards(relPath string) (string, bool) {
	dir, err := os.Getwd()
	if err != nil {
		return "", false
	}

	for {
		candidate := filepath.Join(dir, relPath)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, true
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}
