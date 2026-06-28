package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "spaceship-orders/shared/pkg/proto/inventory/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type InventoryServer struct {
	pb.UnimplementedInventoryServiceServer
	parts map[string]*pb.Part
	mu    sync.RWMutex
}

func (s *InventoryServer) GetPart(ctx context.Context, req *pb.GetPartRequest) (*pb.GetPartResponse, error) {
	part, exists := s.parts[req.Uuid]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "Деталь с UUID %s не найдена", req.Uuid)
	}

	return &pb.GetPartResponse{Part: part}, nil
}

func (s *InventoryServer) ListParts(ctx context.Context, req *pb.ListPartsRequest) (*pb.ListPartsResponse, error) {
	// Блокируем на чтение
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Создаем множество всех возможных тегов для ускорения поиска
	var filterTagsSet map[string]struct{}
	if req.Filter != nil && len(req.Filter.Tags) > 0 {
		filterTagsSet = make(map[string]struct{}, len(req.Filter.Tags))
		for _, tag := range req.Filter.Tags {
			filterTagsSet[tag] = struct{}{}
		}
	}

	// Создаем мапку из всех деталей
	result := make([]*pb.Part, 0, len(s.parts))

	for _, part := range s.parts {

		if req.Filter == nil {
			result = append(result, part)
			continue
		}

		if len(req.Filter.Uuids) > 0 && !containsString(req.Filter.Uuids, part.Uuid) {
			continue
		}

		if len(req.Filter.Names) > 0 && !containsString(req.Filter.Names, part.Name) {
			continue
		}

		if len(req.Filter.Categories) > 0 && !containsCategory(req.Filter.Categories, part.Category) {
			continue
		}

		if len(req.Filter.ManufacturerCountries) > 0 &&
			!containsString(req.Filter.ManufacturerCountries, part.Manufacturer.Country) {
			continue
		}

		if len(req.Filter.Tags) > 0 && !containsTags(req.Filter.Tags, filterTagsSet) {
			continue
		}

		result = append(result, part)

	}

	return &pb.ListPartsResponse{Parts: result}, nil
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func containsCategory(categories []pb.Category, target pb.Category) bool {
	for _, item := range categories {
		if item == target {
			return true
		}
	}
	return false
}

func containsTags(tags []string, target map[string]struct{}) bool {
	for _, tag := range tags {
		if _, exists := target[tag]; exists {
			return true
		}
	}
	return false
}

func main() {
	testPartUUID := "00000000-0000-0000-0000-000000000001"

	srv := &InventoryServer{
		parts: map[string]*pb.Part{
			testPartUUID: {
				Uuid:     testPartUUID,
				Name:     "Тестовый двигатель",
				Price:    500.0,
				Category: pb.Category(1),
				Manufacturer: &pb.Manufacturer{
					Country: "Russia",
				},
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
