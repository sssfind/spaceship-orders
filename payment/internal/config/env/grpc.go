package env

import (
	"errors"
	"net"
	"os"
)

type grpcConfig struct {
	address string
}

func NewPaymentGrpcConfig() (*grpcConfig, error) {
	host := os.Getenv("GRPC_HOST")
	port := os.Getenv("GRPC_PORT")
	if host == "" || port == "" {
		return nil, errors.New("GRPC_HOST or GRPC_PORT is not set")
	}

	return &grpcConfig{
		address: net.JoinHostPort(host, port),
	}, nil
}

func (cfg *grpcConfig) Address() string { return cfg.address }
