package env

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

type orderHTTP struct {
	address string
	timeout time.Duration
}

func NewOrderHTTPConfig() (*orderHTTP, error) {
	host := os.Getenv("HTTP_HOST")
	port := os.Getenv("HTTP_PORT")
	if host == "" || port == "" {
		return nil, errors.New("HHTP_HOST and PORT is not set")
	}

	timeoutStr := os.Getenv("HTTP_READ_TIMEOUT")
	var timeout time.Duration
	var err error

	if timeoutStr == "" {
		timeout = 5 * time.Second
	} else {
		timeout, err = time.ParseDuration(timeoutStr)
		if err != nil {
			return nil, fmt.Errorf("invalid HTTP_READ_TIMEOUT format: %w", err)
		}
	}

	return &orderHTTP{
		address: net.JoinHostPort(host, port),
		timeout: timeout,
	}, nil
}

func (cfg *orderHTTP) GetAddress() string {
	return cfg.address
}

func (cfg *orderHTTP) GetReadTimeout() time.Duration {
	return cfg.timeout
}
