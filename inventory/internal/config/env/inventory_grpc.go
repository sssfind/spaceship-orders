package env

import (
	"errors"
	"net"
	"os"
)

type inventoryGrpcConfig struct {
	address string
}

func NewInventoryGrpcConfig() (*inventoryGrpcConfig, error) {
	host := os.Getenv("GRPC_HOST")
	port := os.Getenv("GRPC_PORT")
	if host == "" || port == "" {
		return nil, errors.New("GRPC_HOST or GRPC_PORT is not set")
	}

	return &inventoryGrpcConfig{
		address: net.JoinHostPort(host, port),
	}, nil
}

func (cfg *inventoryGrpcConfig) Address() string { return cfg.address }
