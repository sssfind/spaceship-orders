package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"order/internal/migrator"
	"os"
	"os/signal"
	"syscall"
	"time"

	apiV1 "order/internal/api/order/v1"
	clientInventory "order/internal/client/grpc/inventory/v1"
	clientPayment "order/internal/client/grpc/payment/v1"
	customMiddleware "order/internal/middleware"
	repoOrder "order/internal/repository/order"
	orderService "order/internal/service/order"
	orderV1 "spaceship-orders/shared/pkg/openapi/order/v1"
	pbInventory "spaceship-orders/shared/pkg/proto/inventory/v1"
	pbPayment "spaceship-orders/shared/pkg/proto/payment/v1"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	httpPort          = "8080"
	readHeaderTimeout = 10 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	initCtx, initCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer initCancel()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/orders?sslmode=disable"
	}

	dbMigrator := migrator.NewMigrator(dsn, "migrations")
	if err := dbMigrator.Up(); err != nil {
		log.Fatalf("migrator up failed: %v", err)
	}
	log.Println("database migrator up")

	dbPool, err := pgxpool.New(initCtx, dsn)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer dbPool.Close()

	invConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("ошибка подключения к InventoryService: %v", err)
	}
	defer func() { _ = invConn.Close() }()
	invGrpcClient := pbInventory.NewInventoryServiceClient(invConn)

	payConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("ошибка подключения к PaymentService: %v", err)
	}
	defer func() { _ = payConn.Close() }()
	payGrpcClient := pbPayment.NewPaymentServiceClient(payConn)

	orderRepo := repoOrder.NewOrderRepository(dbPool)

	inventoryClient := clientInventory.NewClient(invGrpcClient)
	paymentClient := clientPayment.NewClient(payGrpcClient)

	srv := orderService.NewService(orderRepo, inventoryClient, paymentClient)

	apiHandler := apiV1.NewAPI(srv)

	orderServer, err := orderV1.NewServer(apiHandler)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(customMiddleware.RequestLogger)

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("HTTP-сервер запущен на порту %s\n", httpPort)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("Сервер успешно остановлен")
}
