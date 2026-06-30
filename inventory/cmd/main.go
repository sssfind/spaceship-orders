package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	pb "spaceship-orders/shared/pkg/proto/inventory/v1"
)

type InventoryServer struct {
	pb.UnimplementedInventoryServiceServer
	parts map[string]*pb.Part
}

func (s *InventoryServer) GetPart(ctx context.Context, req *pb.GetPartRequest) (*pb.GetPartResponse, error) {
	part, exists := s.parts[req.Uuid]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "Деталь с UUID %s не найдена", req.Uuid)
	}

	return &pb.GetPartResponse{Part: part}, nil
}

func main() {

	srv := &InventoryServer{
		parts: map[string]*pb.Part{
			"1234-5678": {
				Uuid:  "1234-5678",
				Name:  "Тестовый двигатель",
				Price: 100500.50,
			},
		},
	}

	grpcServer := grpc.NewServer()
	pb.RegisterInventoryServiceServer(grpcServer, srv)

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
