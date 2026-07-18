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
	host := os.Getenv("PAYMENT_GRPC_HOST")
	port := os.Getenv("PAYMENT_GRPC_PORT")
	if host == "" || port == "" {
		return nil, errors.New("PAYMENT GRPC_HOST and PORT is not set")
	}

	return &paymentGrpc{
		address: net.JoinHostPort(host, port),
	}, nil
}

func (cfg *paymentGrpc) Address() string { return cfg.address }
