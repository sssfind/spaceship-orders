package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	inventoryV1 "inventory/internal/api/inventory/v1"
	repoPart "inventory/internal/repository/part"
	partService "inventory/internal/service/part"
	pb "spaceship-orders/shared/pkg/proto/inventory/v1"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = godotenv.Load(".env")

	dbURI := os.Getenv("DB_URI")
	if dbURI == "" {
		log.Fatal("DB_URI env var not set")
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051" // Дефолтный порт, если не задан в .env
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Printf("Error connecting to DB %v", err)
		return
	}

	defer func() {
		disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer disconnectCancel()
		if err := client.Disconnect(disconnectCtx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("MongoDB недоступна, ошибка ping: %v\n", err)
		return
	}
	log.Println("Успешное подключение к MongoDB")

	collection := client.Database("spaceship-orders").Collection("parts")

	repo := repoPart.NewPartRepository(collection)
	srv := partService.NewService(repo)
	apiHandler := inventoryV1.NewAPI(srv)

	grpcServer := grpc.NewServer()
	pb.RegisterInventoryServiceServer(grpcServer, apiHandler)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", net.JoinHostPort("0.0.0.0", grpcPort))
	if err != nil {
		log.Fatalf("Ошибка при прослушивании порта %s: %v", grpcPort, err)
	}

	go func() {
		fmt.Printf("InventoryService запущен на порту %s\n", grpcPort)
		if err := grpcServer.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			log.Fatalf("Ошибка gRPC сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Завершение работы InventoryService...")

	grpcServer.GracefulStop()
	log.Println("InventoryService успешно остановлен")
}
