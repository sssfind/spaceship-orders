package env

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

type httpConfig struct {
	address string
	timeout time.Duration
}

func NewHTTPConfig() (*httpConfig, error) {
	host := os.Getenv("HTTP_HOST")
	port := os.Getenv("HTTP_PORT")
	if host == "" || port == "" {
		return nil, errors.New("HTTP_HOST and HTTP_PORT are not set")
	}

	timeoutStr := os.Getenv("HTTP_READ_TIMEOUT")
	timeout := 5 * time.Second
	if timeoutStr != "" {
		var err error
		timeout, err = time.ParseDuration(timeoutStr)
		if err != nil {
			return nil, fmt.Errorf("invalid HTTP_READ_TIMEOUT format: %w", err)
		}
	}

	return &httpConfig{
		address: net.JoinHostPort(host, port),
		timeout: timeout,
	}, nil
}

func (cfg *httpConfig) Address() string            { return cfg.address }
func (cfg *httpConfig) ReadTimeout() time.Duration { return cfg.timeout }
