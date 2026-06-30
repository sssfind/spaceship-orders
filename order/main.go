package main

import (
	"context"
	"log"
	orderV1 "spaceship-orders/shared/pkg/openapi/order/v1"
	pbInventory "spaceship-orders/shared/pkg/proto/inventory/v1"
	pbPayment "spaceship-orders/shared/pkg/proto/payment/v1"

	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	httpPort        = "8080"
	readerTimeout   = 10 * time.Second
	shutdownTimeout = 10 * time.Second
)

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.Order),
	}
}

func (s *OrderStorage) Get(uuid string) *orderV1.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ok := s.orders[uuid]
	if !ok {
		return nil
	}
	return order
}

func (s *OrderStorage) Create(order *orderV1.Order) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[order.OrderUUID.String()] = order
}

func (s *OrderStorage) UpdateStatus(orderID string, status string, txUUID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	order, ok := s.orders[orderID]
	if !ok {
		return
	}

	if status != "" {
		order.Status = orderV1.OrderStatus(status)
	}

	if txUUID != "" {
		parsedUUID, err := uuid.Parse(txUUID)
		if err != nil {
			return
		}

		order.TransactionUUID = orderV1.NewOptNilUUID(parsedUUID)
	}
}

type OrderHandler struct {
	storage   *OrderStorage
	invClient pbInventory.InventoryServiceClient
	payClient pbPayment.PaymentServiceClient
}

func NewOrderHandler(storage *OrderStorage, inv pbInventory.InventoryServiceClient, pay pbPayment.PaymentServiceClient) *OrderHandler {
	return &OrderHandler{
		storage:   storage,
		invClient: inv,
		payClient: pay,
	}
}

func (h *OrderHandler) GetOrderByUUID(_ context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	order := h.storage.Get(params.OrderUUID.String())
	if order == nil {
		return &orderV1.GetOrderByUUIDNotFound{}, nil
	}
	return order, nil
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderReq) (orderV1.CreateOrderRes, error) {
	partUUIDsStrings := make([]string, 0, len(req.PartUuids))

	for _, uuid := range req.PartUuids {
		partUUIDsStrings = append(partUUIDsStrings, uuid.String())
	}

	inventoryReq := &pbInventory.ListPartsRequest{
		Filter: &pbInventory.PartsFilter{
			Uuids: partUUIDsStrings,
		},
	}

	h.invClient.ListParts(ctx, inventoryReq)
}

func main() {
	// Подключаемся к инвентарю
	invConn, _ := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	invClient := pbInventory.NewInventoryServiceClient(invConn)

	// Подключаемся к оплате
	payConn, _ := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	payClient := pbPayment.NewPaymentServiceClient(payConn)

	storage := NewOrderStorage()

	handler := NewOrderHandler(storage, invClient, payClient)

	orderServer, err := orderV1.NewServer(handler)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(customMiddleware.RequestLogger)

	r.Mount("/", orderServer)
}
