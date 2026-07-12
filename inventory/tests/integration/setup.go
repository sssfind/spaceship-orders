//go:build integration

package integration

import (
	"context"
	"net"
	"os"
	"time"

	"inventory/internal/app"
	"spaceship-orders/shared/pkg/proto/inventory/v1"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (s *InventoryTestSuite) SetupSuite() {
	s.ctx = context.Background()

	mongoReq := testcontainers.ContainerRequest{
		Image:        "mongo:6.0",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForLog("Waiting for connections"),
	}

	container, err := testcontainers.GenericContainer(s.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: mongoReq,
		Started:          true,
	})
	s.Require().NoError(err)
	s.mongoContainer = container

	mongoHost, err := container.Host(s.ctx)
	s.Require().NoError(err)
	mongoPort, err := container.MappedPort(s.ctx, "27017")
	s.Require().NoError(err)

	testPort := "50091"
	os.Setenv("GRPC_HOST", "127.0.0.1")
	os.Setenv("GRPC_PORT", testPort)
	os.Setenv("MONGO_HOST", mongoHost)
	os.Setenv("MONGO_PORT", mongoPort.Port())
	os.Setenv("MONGO_DATABASE", "inventory_integration_test")
	os.Setenv("LOGGER_LEVEL", "fatal")

	application, err := app.NewApp(s.ctx)
	s.Require().NoError(err)
	s.app = application

	go func() {
		_ = s.app.Run()
	}()
	time.Sleep(100 * time.Millisecond)

	conn, err := grpc.NewClient(
		net.JoinHostPort("127.0.0.1", testPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.Require().NoError(err)
	s.grpcConn = conn
	s.client = v1.NewInventoryServiceClient(conn)
}
