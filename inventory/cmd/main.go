package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	inventoryV1 "inventory/internal/api/inventory/v1"
	repoPart "inventory/internal/repository/part"
	partService "inventory/internal/service/part"
	pb "spaceship-orders/shared/pkg/proto/inventory/v1"
)

func main() {
	repo := repoPart.NewPartRepository()
	srv := partService.NewService(repo)
	apiHandler := inventoryV1.NewAPI(srv)

	grpcServer := grpc.NewServer()
	pb.RegisterInventoryServiceServer(grpcServer, apiHandler)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Ошибка при прослушивании порта: %v", err)
	}

	fmt.Println("🚀 InventoryService запущен на порту 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка сервера: %v", err)
	}
}
