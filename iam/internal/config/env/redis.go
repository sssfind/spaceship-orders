package env

import (
	"errors"
	"net"
	"os"
	"strconv"
	"time"
)

type redisConfig struct {
	host              string
	port              string
	connectionTimeout time.Duration
	maxIdle           int
	idleTimeout       time.Duration
}

func NewRedisConfig() (*redisConfig, error) {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		return nil, errors.New("REDIS_HOST is not set")
	}
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		return nil, errors.New("REDIS_PORT is not set")
	}

	connTimeoutStr := os.Getenv("REDIS_CONNECTION_TIMEOUT")
	connTimeout, err := time.ParseDuration(connTimeoutStr)
	if err != nil {
		connTimeout = 5 * time.Second
	}

	maxIdleStr := os.Getenv("REDIS_MAX_IDLE")
	maxIdle, err := strconv.Atoi(maxIdleStr)
	if err != nil {
		maxIdle = 10
	}

	idleTimeoutStr := os.Getenv("REDIS_IDLE_TIMEOUT")
	idleTimeout, err := time.ParseDuration(idleTimeoutStr)
	if err != nil {
		idleTimeout = 5 * time.Minute
	}

	return &redisConfig{
		host:              host,
		port:              port,
		connectionTimeout: connTimeout,
		maxIdle:           maxIdle,
		idleTimeout:       idleTimeout,
	}, nil
}

func (r *redisConfig) Address() string {
	return net.JoinHostPort(r.host, r.port)
}

func (r *redisConfig) ConnectionTimeout() time.Duration {
	return r.connectionTimeout
}

func (r *redisConfig) MaxIdle() int {
	return r.maxIdle
}

func (r *redisConfig) IdleTimeout() time.Duration {
	return r.idleTimeout
}
