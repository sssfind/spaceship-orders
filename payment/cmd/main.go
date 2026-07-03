package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	paymentV1 "payment/internal/api/payment/v1"
	paymentService "payment/internal/service/payment"
	pb "spaceship-orders/shared/pkg/proto/payment/v1"
)

const grpcPort = 50052

func main() {
	srv := paymentService.NewService()
	apiHandler := paymentV1.NewAPI(srv)

	grpcServer := grpc.NewServer()
	pb.RegisterPaymentServiceServer(grpcServer, apiHandler)

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("Ошибка при прослушивании порта: %v", err)
	}

	fmt.Printf("PaymentService запущен на порту %d", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка сервера: %v", err)
	}
}
