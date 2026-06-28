package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "spaceship-orders/shared/pkg/proto/payment/v1"
)

const grpcPort = 50052

type PaymentServiceServer struct {
	pb.UnimplementedPaymentServiceServer
}

func (s *PaymentServiceServer) PayOrder(ctx context.Context, req *pb.PayOrderRequest) (*pb.PayOrderResponse, error) {
	payUuid := uuid.NewString()
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", payUuid)
	return &pb.PayOrderResponse{TransactionUuid: payUuid}, nil
}

func main() {
	srv := &PaymentServiceServer{}

	grpcServer := grpc.NewServer()
	pb.RegisterPaymentServiceServer(grpcServer, srv)

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
