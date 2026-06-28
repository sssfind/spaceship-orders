package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	customMiddleware "order/internal/middleware"
	orderV1 "spaceship-orders/shared/pkg/openapi/order/v1"
	pbInventory "spaceship-orders/shared/pkg/proto/inventory/v1"
	pbPayment "spaceship-orders/shared/pkg/proto/payment/v1"
)

const (
	httpPort          = "8080"
	readHeaderTimeout = 10 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.OrderDto
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.OrderDto),
	}
}

func (s *OrderStorage) Get(uuid string) *orderV1.OrderDto {
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ok := s.orders[uuid]
	if !ok {
		return nil
	}
	return order
}

func (s *OrderStorage) Create(order *orderV1.OrderDto) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[order.OrderUUID.String()] = order
}

func (s *OrderStorage) UpdateStatus(orderID string, status orderV1.OrderStatus, txUUID string, payMethod orderV1.PaymentMethod) {
	s.mu.Lock()
	defer s.mu.Unlock()
	order, ok := s.orders[orderID]
	if !ok {
		return
	}

	if status != "" {
		order.Status = status
	}

	if txUUID != "" {
		if parsedUUID, err := uuid.Parse(txUUID); err == nil {
			order.TransactionUUID = orderV1.NewOptUUID(parsedUUID)
		}
	}

	if payMethod != "" {
		order.PaymentMethod = orderV1.NewOptNilPaymentMethod(payMethod)
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
		return &orderV1.GenericError{
			Code:    404,
			Message: "Order not found",
		}, nil
	}
	return order, nil
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	partUUIDsStrings := make([]string, 0, len(req.PartUuids))

	for _, pUUID := range req.PartUuids {
		partUUIDsStrings = append(partUUIDsStrings, pUUID.String())
	}

	inventoryReq := &pbInventory.ListPartsRequest{
		Filter: &pbInventory.PartsFilter{
			Uuids: partUUIDsStrings,
		},
	}

	inventoryRes, err := h.invClient.ListParts(ctx, inventoryReq)
	if err != nil {
		return &orderV1.CreateOrderInternalServerError{
			Code:    500,
			Message: "Inventory service is unavailable: " + err.Error(),
		}, nil
	}

	if len(inventoryRes.Parts) != len(req.PartUuids) {
		return &orderV1.CreateOrderBadRequest{
			Code:    400,
			Message: "Some requested parts were not found in inventory",
		}, nil
	}

	userUUID, err := uuid.Parse(req.UserUUID)
	if err != nil {
		return &orderV1.CreateOrderBadRequest{
			Code:    400,
			Message: "Invalid user_uuid format",
		}, nil
	}

	var totalPrice float64
	for _, part := range inventoryRes.Parts {
		totalPrice += part.Price
	}

	orderUUID := uuid.New()
	newOrder := &orderV1.OrderDto{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUuids:  req.PartUuids,
		TotalPrice: totalPrice,
		Status:     orderV1.OrderStatusPENDINGPAYMENT,
	}
	h.storage.Create(newOrder)

	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}

func (h *OrderHandler) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	order := h.storage.Get(params.OrderUUID.String())
	if order == nil {
		return &orderV1.CancelOrderNotFound{
			Code:    404,
			Message: "Order not exists",
		}, nil
	}

	if order.Status == orderV1.OrderStatusPAID {
		return &orderV1.CancelOrderConflict{
			Code:    409,
			Message: "The order has already been paid",
		}, nil
	}

	if order.Status == orderV1.OrderStatusPENDINGPAYMENT {
		h.storage.UpdateStatus(order.OrderUUID.String(), orderV1.OrderStatusCANCELLED, "", "")
	}

	return &orderV1.CancelOrderNoContent{}, nil
}

func (h *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	order := h.storage.Get(params.OrderUUID.String())
	if order == nil {
		return &orderV1.PayOrderNotFound{
			Code:    404,
			Message: "Order not found",
		}, nil
	}

	if order.Status != orderV1.OrderStatusPENDINGPAYMENT {
		return &orderV1.PayOrderBadRequest{
			Code:    400,
			Message: "Invalid order status for payment",
		}, nil
	}

	var grpcPaymentMethod pbPayment.PaymentMethod
	switch req.PaymentMethod {
	case orderV1.PaymentMethodCARD:
		grpcPaymentMethod = pbPayment.PaymentMethod_PAYMENT_METHOD_CARD
	case orderV1.PaymentMethodSBP:
		grpcPaymentMethod = pbPayment.PaymentMethod_PAYMENT_METHOD_SBP
	case orderV1.PaymentMethodCREDITCARD:
		grpcPaymentMethod = pbPayment.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderV1.PaymentMethodINVESTORMONEY:
		grpcPaymentMethod = pbPayment.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		grpcPaymentMethod = pbPayment.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}

	payRes, err := h.payClient.PayOrder(ctx, &pbPayment.PayOrderRequest{
		OrderUuid:     order.OrderUUID.String(),
		UserUuid:      order.UserUUID.String(),
		PaymentMethod: grpcPaymentMethod,
	})
	if err != nil {
		return &orderV1.PayOrderInternalServerError{
			Code:    500,
			Message: "Payment processing failed: " + err.Error(),
		}, nil
	}

	h.storage.UpdateStatus(order.OrderUUID.String(), orderV1.OrderStatusPAID, payRes.TransactionUuid, req.PaymentMethod)

	txUUID, _ := uuid.Parse(payRes.TransactionUuid)
	return &orderV1.PayOrderResponse{
		TransactionUUID: txUUID,
	}, nil
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

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("Сервер остановлен")
}
