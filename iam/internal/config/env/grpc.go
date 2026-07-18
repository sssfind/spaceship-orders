package env

import (
	"errors"
	"net"
	"os"
)

type grpcConfig struct {
	host string
	port string
}

func NewGRPCConfig() (*grpcConfig, error) {
	host := os.Getenv("GRPC_HOST")
	if host == "" {
		return nil, errors.New("GRPC_HOST is not set")
	}
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		return nil, errors.New("GRPC_PORT is not set")
	}

	return &grpcConfig{host: host, port: port}, nil
}

func (g *grpcConfig) Address() string {
	return net.JoinHostPort(g.host, g.port)
}
