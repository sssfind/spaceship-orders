package part

import (
	"inventory/internal/repository"

	"go.mongodb.org/mongo-driver/mongo"
)

type repo struct {
	collection *mongo.Collection
}

func NewPartRepository(collection *mongo.Collection) repository.PartRepository {
	r := &repo{
		collection: collection,
	}
	r.initTestData() // Заполняем тестовыми данными из init.go
	return r
}
