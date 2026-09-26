package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"order/internal/config/env"
)

type Config struct {
	LoggerConfig
	TracerConfig
	OrderHttpConfig
	PostgresConfig
	IAMHTTPConfig
}

func Load() (*Config, error) {
	loadDotEnv()

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return nil, fmt.Errorf("logger config error: %w", err)
	}

	tracerCfg, err := env.NewTracerConfig()
	if err != nil {
		return nil, fmt.Errorf("tracer config error: %w", err)
	}

	httpCfg, err := env.NewOrderHTTPConfig()
	if err != nil {
		return nil, fmt.Errorf("http config error: %w", err)
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return nil, fmt.Errorf("postgres config error: %w", err)
	}

	iamCfg, err := env.NewIAMHTTPConfig()
	if err != nil {
		return nil, fmt.Errorf("iam http config error: %w", err)
	}

	return &Config{
		LoggerConfig:    loggerCfg,
		TracerConfig:    tracerCfg,
		OrderHttpConfig: httpCfg,
		PostgresConfig:  postgresCfg,
		IAMHTTPConfig:   iamCfg,
	}, nil
}

func (c *Config) ServiceName() string {
	return c.LoggerConfig.ServiceName()
}

func loadDotEnv() {
	candidates := []string{".env"}

	if composeEnv, ok := findFileUpwards("deploy/compose/order/.env"); ok {
		candidates = append(candidates, composeEnv)
	} else {
		candidates = append(candidates,
			filepath.Join("..", "deploy", "compose", "order", ".env"),
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
