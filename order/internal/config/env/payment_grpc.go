package env

import (
	"errors"
	"net"
	"os"
)

type paymentGrpc struct {
	address string
}

func NewPaymentGrpcConfig() (*paymentGrpc, error) {
	host := os.Getenv("INVENTORY_GRPC_HOST")
	port := os.Getenv("INVENTORY_GRPC_PORT")
	if host == "" || port == "" {
		return nil, errors.New("GRPC_HOST and PORT is not set")
	}

	return &paymentGrpc{
		address: net.JoinHostPort(host, port),
	}, nil
}

func (cfg *paymentGrpc) Address() string { return cfg.address }
