package part

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	repoModel "inventory/internal/repository/model"
)

func (r *repo) initTestData() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testPartUUID := "00000000-0000-0000-0000-000000000001"

	part := &repoModel.Part{
		UUID:     testPartUUID,
		Name:     "Тестовый двигатель",
		Price:    500.0,
		Category: repoModel.Category(1),
	}

	filter := bson.M{"uuid": testPartUUID}
	update := bson.M{"$set": part}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Printf("inventory repository initTestData error: %v", err)
	}
}
