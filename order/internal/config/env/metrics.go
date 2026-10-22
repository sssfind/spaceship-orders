package env

import (
	"net"
	"os"
)

type metricsHTTP struct {
	address string
}

func NewMetricsHTTPConfig() (*metricsHTTP, error) {
	host := os.Getenv("METRICS_HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	port := os.Getenv("METRICS_PORT")
	if port == "" {
		port = "9090"
	}

	return &metricsHTTP{
		address: net.JoinHostPort(host, port),
	}, nil
}

func (cfg *metricsHTTP) Address() string {
	return cfg.address
}

