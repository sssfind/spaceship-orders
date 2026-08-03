package env

import (
	"net"
	"os"
)

const (
	httpHostEnv = "ASSEMBLY_HTTP_HOST"
	httpPortEnv = "ASSEMBLY_HTTP_PORT"
)

type HTTPConfig interface {
	Address() string
}

type httpConfig struct {
	host string
	port string
}

func NewHTTPConfig() (HTTPConfig, error) {
	host := os.Getenv(httpHostEnv)
	if len(host) == 0 {
		host = "0.0.0.0"
	}

	port := os.Getenv(httpPortEnv)
	if len(port) == 0 {
		port = "8082"
	}

	return &httpConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *httpConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
