package env

import (
	"errors"
	"net"
	"os"
)

type inventoryGrpc struct {
	address string
}

func NewInventoryGrpcConfig() (*inventoryGrpc, error) {
	host := os.Getenv("INVENTORY_GRPC_HOST")
	port := os.Getenv("INVENTORY_GRPC_PORT")
	if host == "" || port == "" {
		return nil, errors.New("GRPC_HOST and PORT is not set")
	}

	return &inventoryGrpc{
		address: net.JoinHostPort(host, port),
	}, nil
}

func (cfg *inventoryGrpc) Address() string { return cfg.address }
