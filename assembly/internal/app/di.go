package app

import (
	"assembly/internal/config"
	"assembly/internal/service/consumer/order_consumer"
	"assembly/internal/service/producer/order_producer"
	"github.com/IBM/sarama"
	"platform/pkg/kafka/consumer"
	"platform/pkg/kafka/producer"
	"platform/pkg/logger"
)

type serviceProvider struct {
	cfg config.Config

	saramaProducer sarama.SyncProducer
	saramaConsumer sarama.ConsumerGroup

	assemblyProducer order_producer.AssemblyProducer
	orderConsumer    *order_consumer.OrderConsumer
}

func newServiceProvider(cfg config.Config) *serviceProvider {
	return &serviceProvider{
		cfg: cfg,
	}
}

func (sp *serviceProvider) SaramaProducer() (sarama.SyncProducer, error) {
	if sp.saramaProducer == nil {
		saramaCfg := sarama.NewConfig()
		saramaCfg.Producer.RequiredAcks = sarama.WaitForAll
		saramaCfg.Producer.Return.Successes = true

		prod, err := sarama.NewSyncProducer(sp.cfg.Brokers(), saramaCfg)
		if err != nil {
			return nil, err
		}
		sp.saramaProducer = prod
	}
	return sp.saramaProducer, nil
}

func (sp *serviceProvider) SaramaConsumerGroup() (sarama.ConsumerGroup, error) {
	if sp.saramaConsumer == nil {
		saramaCfg := sarama.NewConfig()
		saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest

		group, err := sarama.NewConsumerGroup(sp.cfg.Brokers(), sp.cfg.GroupID(), saramaCfg)
		if err != nil {
			return nil, err
		}
		sp.saramaConsumer = group
	}
	return sp.saramaConsumer, nil
}

func (sp *serviceProvider) AssemblyProducer() (order_producer.AssemblyProducer, error) {
	if sp.assemblyProducer == nil {
		syncProd, err := sp.SaramaProducer()
		if err != nil {
			return nil, err
		}

		platformProd := producer.NewProducer(syncProd, sp.cfg.AssembledTopic(), logger.Logger())
		sp.assemblyProducer = order_producer.NewAssemblyProducer(platformProd)
	}
	return sp.assemblyProducer, nil
}

func (sp *serviceProvider) OrderConsumer() (*order_consumer.OrderConsumer, error) {
	if sp.orderConsumer == nil {
		group, err := sp.SaramaConsumerGroup()
		if err != nil {
			return nil, err
		}

		prod, err := sp.AssemblyProducer()
		if err != nil {
			return nil, err
		}

		// Передаем глобальный инициализированный логгер
		platformCons := consumer.NewConsumer(group, []string{sp.cfg.PaidTopic()}, logger.Logger())

		handler := order_consumer.NewOrderPaidHandler(prod)
		sp.orderConsumer = order_consumer.NewOrderConsumer(platformCons, handler)
	}
	return sp.orderConsumer, nil
}
