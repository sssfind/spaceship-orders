package part

import (
	"context"
	"errors"
	"fmt"
	repoModel "inventory/internal/repository/model"

	"inventory/internal/model"
	"inventory/internal/repository/converter"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Реализация метода Get для интерфейса
func (r *repo) Get(ctx context.Context, uuid string) (*model.Part, error) {
	filter := bson.M{"uuid": uuid}

	var part repoModel.Part
	err := r.collection.FindOne(ctx, filter).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, model.ErrPartNotFound
		}
		return nil, fmt.Errorf("get part: %w", err)
	}
	domainPart := converter.PartToDomain(&part)
	return &domainPart, nil
}
