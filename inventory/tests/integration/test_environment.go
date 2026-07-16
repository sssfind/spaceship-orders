//go:build integration

package integration

import (
	"context"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"google.golang.org/grpc"
	"inventory/internal/app"
	"spaceship-orders/shared/pkg/proto/inventory/v1"
)

type InventoryTestSuite struct {
	suite.Suite
	ctx            context.Context
	mongoContainer testcontainers.Container
	app            *app.App
	grpcConn       *grpc.ClientConn
	client         v1.InventoryServiceClient
}
