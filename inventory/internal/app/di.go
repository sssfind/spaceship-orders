package app

import (
	"context"
	"fmt"
	"inventory/internal/config"
	"inventory/internal/repository"
	repoPart "inventory/internal/repository/part"
	"inventory/internal/service"
	implPart "inventory/internal/service/part"
	"platform/pkg/closer"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type serviceProvider struct {
	cfg         *config.Config
	mongoClient *mongo.Client
	mongoDB     *mongo.Database
	partRepo    repository.PartRepository
	partSrv     service.PartService
}

func newServiceProvider(cfg *config.Config) *serviceProvider {
	return &serviceProvider{cfg: cfg}
}

// MongoClient инициализирует пул соединений с MongoDB
func (sp *serviceProvider) MongoClient(ctx context.Context) (*mongo.Client, error) {
	if sp.mongoClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(sp.cfg.GetURI()))
		if err != nil {
			return nil, fmt.Errorf("failed to connect to mongo: %w", err)
		}

		if err := client.Ping(ctx, nil); err != nil {
			return nil, fmt.Errorf("failed to ping mongo: %w", err)
		}

		sp.mongoClient = client

		closer.AddNamed("mongodb_client", func(c context.Context) error {
			return sp.mongoClient.Disconnect(c)
		})
	}
	return sp.mongoClient, nil
}

func (sp *serviceProvider) MongoDatabase(ctx context.Context) (*mongo.Database, error) {
	if sp.mongoDB == nil {
		client, err := sp.MongoClient(ctx)
		if err != nil {
			return nil, err
		}
		sp.mongoDB = client.Database(sp.cfg.GetDatabaseName())
	}
	return sp.mongoDB, nil
}

func (sp *serviceProvider) PartRepository(ctx context.Context) (repository.PartRepository, error) {
	if sp.partRepo == nil {
		db, err := sp.MongoDatabase(ctx)
		if err != nil {
			return nil, err
		}
		sp.partRepo = repoPart.NewPartRepository(db.Collection("parts"))
	}
	return sp.partRepo, nil
}

func (sp *serviceProvider) PartService(ctx context.Context) (service.PartService, error) {
	if sp.partSrv == nil {
		repo, err := sp.PartRepository(ctx)
		if err != nil {
			return nil, err
		}
		sp.partSrv = implPart.NewService(repo)
	}
	return sp.partSrv, nil
}
