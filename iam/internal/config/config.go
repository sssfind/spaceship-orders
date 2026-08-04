package config

import (
	"github.com/joho/godotenv"
	"iam/internal/config/env"
)

type Config struct {
	grpc     GRPCConfig
	logger   LoggerConfig
	postgres PostgresConfig
	redis    RedisConfig
	session  SessionConfig
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	grpcCfg, err := env.NewGRPCConfig()
	if err != nil {
		return nil, err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return nil, err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return nil, err
	}

	redisCfg, err := env.NewRedisConfig()
	if err != nil {
		return nil, err
	}

	sessionCfg, err := env.NewSessionConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		grpc:     grpcCfg,
		logger:   loggerCfg,
		postgres: postgresCfg,
		redis:    redisCfg,
		session:  sessionCfg,
	}, nil
}

func (c *Config) GRPC() GRPCConfig         { return c.grpc }
func (c *Config) Logger() LoggerConfig     { return c.logger }
func (c *Config) Postgres() PostgresConfig { return c.postgres }
func (c *Config) Redis() RedisConfig       { return c.redis }
func (c *Config) Session() SessionConfig   { return c.session }
