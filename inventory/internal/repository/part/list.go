package part

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	domainModel "inventory/internal/model"
	"inventory/internal/repository/converter"
	repoModel "inventory/internal/repository/model"
)

func (r *repo) List(ctx context.Context, filter *domainModel.PartsFilter) ([]domainModel.Part, error) {
	mongoFilter := bson.M{}

	var filterTagsSet map[string]struct{}
	if filter != nil && len(filter.Tags) > 0 {
		filterTagsSet = make(map[string]struct{}, len(filter.Tags))
		for _, tag := range filter.Tags {
			filterTagsSet[tag] = struct{}{}
		}
	}

	if filter != nil {
		if len(filter.UUIDs) > 0 {
			mongoFilter["uuid"] = bson.M{"$in": filter.UUIDs}
		}
		if len(filter.Names) > 0 {
			mongoFilter["name"] = bson.M{"$in": filter.Names}
		}
		if len(filter.Categories) > 0 {
			mongoFilter["category"] = bson.M{"$in": filter.Categories}
		}
		if len(filter.ManufacturerCountries) > 0 {
			mongoFilter["manufacturer.country"] = bson.M{"$in": filter.ManufacturerCountries}
		}
		if len(filter.Tags) > 0 {
			// В Mongo оператор $in для массива поля означает: "выбрать документы,
			// где в массиве tags есть ХОТЯ БЫ ОДИН элемент из filter.Tags"
			mongoFilter["tags"] = bson.M{"$in": filter.Tags}
		}
	}

	cursor, err := r.collection.Find(ctx, mongoFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to list parts: %w", err)
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("failed to close cursor: %v", err)
		}
	}(cursor, ctx)

	var result []domainModel.Part

	for cursor.Next(ctx) {
		var part repoModel.Part
		if err := cursor.Decode(&part); err != nil {
			return nil, fmt.Errorf("failed to decode part: %w", err)
		}
		result = append(result, converter.PartToDomain(&part))
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor err: %w", err)
	}

	return result, nil
}
