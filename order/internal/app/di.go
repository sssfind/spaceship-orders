package app

import (
	"context"
	"fmt"

	apiV1 "order/internal/api/order/v1"
	clientInventory "order/internal/client/grpc/inventory/v1"
	clientPayment "order/internal/client/grpc/payment/v1"
	"order/internal/config"
	"order/internal/consumer/order_consumer"
	"order/internal/producer/order_producer"
	"order/internal/repository"
	repoOrder "order/internal/repository/order"
	"order/internal/service"
	orderImpl "order/internal/service/order"
	"platform/pkg/closer"
	platformConsumer "platform/pkg/kafka/consumer"
	platformProducer "platform/pkg/kafka/producer"
	"platform/pkg/logger"
	orderV1 "spaceship-orders/shared/pkg/openapi/order/v1"
	pbInventory "spaceship-orders/shared/pkg/proto/inventory/v1"
	pbPayment "spaceship-orders/shared/pkg/proto/payment/v1"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type serviceProvider struct {
	cfg *config.Config

	dbPool          *pgxpool.Pool
	invConn         *grpc.ClientConn
	payConn         *grpc.ClientConn
	inventoryClient orderImpl.InventoryClient
	paymentClient   orderImpl.PaymentClient
	orderRepo       repository.OrderRepository
	orderSrv        service.OrderService
	apiHandler      orderV1.Handler

	saramaProducer sarama.SyncProducer
	saramaConsumer sarama.ConsumerGroup
	orderProducer  order_producer.OrderProducer
	orderConsumer  *order_consumer.OrderConsumer
}

func newServiceProvider(cfg *config.Config) *serviceProvider {
	return &serviceProvider{cfg: cfg}
}

func (sp *serviceProvider) DBPool(ctx context.Context) (*pgxpool.Pool, error) {
	if sp.dbPool == nil {
		pool, err := pgxpool.New(ctx, sp.cfg.Dsn())
		if err != nil {
			return nil, fmt.Errorf("failed to connect to postgres: %w", err)
		}
		sp.dbPool = pool

		closer.AddNamed("postgres_pool", func(_ context.Context) error {
			sp.dbPool.Close()
			return nil
		})
	}
	return sp.dbPool, nil
}

func (sp *serviceProvider) InventoryGrpcClient(ctx context.Context) (orderImpl.InventoryClient, error) {
	if sp.inventoryClient == nil {
		conn, err := grpc.NewClient(sp.cfg.InventoryGrpcConfig.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("inventory connection error: %w", err)
		}
		sp.invConn = conn

		closer.AddNamed("grpc_inventory_connection", func(_ context.Context) error {
			return sp.invConn.Close()
		})

		pbClient := pbInventory.NewInventoryServiceClient(sp.invConn)
		sp.inventoryClient = clientInventory.NewClient(pbClient)
	}
	return sp.inventoryClient, nil
}

func (sp *serviceProvider) PaymentGrpcClient(ctx context.Context) (orderImpl.PaymentClient, error) {
	if sp.paymentClient == nil {
		conn, err := grpc.NewClient(sp.cfg.PaymentGrpcConfig.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("payment connection error: %w", err)
		}
		sp.payConn = conn

		closer.AddNamed("grpc_payment_connection", func(_ context.Context) error {
			return sp.payConn.Close()
		})

		pbClient := pbPayment.NewPaymentServiceClient(sp.payConn)
		sp.paymentClient = clientPayment.NewClient(pbClient)
	}
	return sp.paymentClient, nil
}

func (sp *serviceProvider) OrderRepository(ctx context.Context) (repository.OrderRepository, error) {
	if sp.orderRepo == nil {
		pool, err := sp.DBPool(ctx)
		if err != nil {
			return nil, err
		}
		sp.orderRepo = repoOrder.NewOrderRepository(pool)
	}
	return sp.orderRepo, nil
}

func (sp *serviceProvider) SaramaProducer() (sarama.SyncProducer, error) {
	if sp.saramaProducer == nil {
		saramaCfg := sarama.NewConfig()
		saramaCfg.Producer.RequiredAcks = sarama.WaitForAll
		saramaCfg.Producer.Return.Successes = true

		prod, err := sarama.NewSyncProducer(sp.cfg.Brokers(), saramaCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create sarama producer: %w", err)
		}
		sp.saramaProducer = prod

		closer.AddNamed("order_sarama_producer", func(_ context.Context) error {
			return sp.saramaProducer.Close()
		})
	}
	return sp.saramaProducer, nil
}

func (sp *serviceProvider) SaramaConsumerGroup() (sarama.ConsumerGroup, error) {
	if sp.saramaConsumer == nil {
		saramaCfg := sarama.NewConfig()
		saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest

		group, err := sarama.NewConsumerGroup(sp.cfg.Brokers(), sp.cfg.GroupID(), saramaCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create sarama consumer group: %w", err)
		}
		sp.saramaConsumer = group

		closer.AddNamed("order_sarama_consumer_group", func(_ context.Context) error {
			return sp.saramaConsumer.Close()
		})
	}
	return sp.saramaConsumer, nil
}

func (sp *serviceProvider) OrderProducer() (order_producer.OrderProducer, error) {
	if sp.orderProducer == nil {
		syncProd, err := sp.SaramaProducer()
		if err != nil {
			return nil, err
		}

		platformProd := platformProducer.NewProducer(syncProd, sp.cfg.PaidTopic(), logger.Logger())
		sp.orderProducer = order_producer.NewOrderProducer(platformProd)
	}
	return sp.orderProducer, nil
}

func (sp *serviceProvider) OrderConsumer(ctx context.Context) (*order_consumer.OrderConsumer, error) {
	if sp.orderConsumer == nil {
		group, err := sp.SaramaConsumerGroup()
		if err != nil {
			return nil, err
		}

		orderSrv, err := sp.OrderService(ctx)
		if err != nil {
			return nil, err
		}

		platformCons := platformConsumer.NewConsumer(group, []string{sp.cfg.AssembledTopic()}, logger.Logger())

		handler := order_consumer.NewOrderAssembledHandler(orderSrv)
		sp.orderConsumer = order_consumer.NewOrderConsumer(platformCons, handler)
	}
	return sp.orderConsumer, nil
}

func (sp *serviceProvider) OrderService(ctx context.Context) (service.OrderService, error) {
	if sp.orderSrv == nil {
		repo, err := sp.OrderRepository(ctx)
		if err != nil {
			return nil, err
		}
		invClient, err := sp.InventoryGrpcClient(ctx)
		if err != nil {
			return nil, err
		}
		payClient, err := sp.PaymentGrpcClient(ctx)
		if err != nil {
			return nil, err
		}
		prod, err := sp.OrderProducer()
		if err != nil {
			return nil, err
		}

		sp.orderSrv = orderImpl.NewService(repo, invClient, payClient, prod)
	}
	return sp.orderSrv, nil
}

func (sp *serviceProvider) APIHandler(ctx context.Context) (orderV1.Handler, error) {
	if sp.apiHandler == nil {
		srv, err := sp.OrderService(ctx)
		if err != nil {
			return nil, err
		}
		sp.apiHandler = apiV1.NewAPI(srv)
	}
	return sp.apiHandler, nil
}
