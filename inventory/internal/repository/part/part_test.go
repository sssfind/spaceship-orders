package part_test

import (
	"context"
	"testing"
	"time"

	"inventory/internal/model"
	"inventory/internal/repository/part"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestCollection(t *testing.T) *mongo.Collection {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	mongoURI := "mongodb://root:secretpassword@localhost:27017/?authSource=admin"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("ошибка подключения к тестовой Mongo: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		t.Fatalf("ошибка подключения к тестовой Mongo (сервер недоступен): %v", err)
	}

	db := client.Database("spaceship-orders-test")
	collection := db.Collection("parts")

	if _, err := collection.DeleteMany(ctx, bson.M{}); err != nil {
		t.Fatalf("ошибка очистки тестовой коллекции: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Drop(context.Background())
		_ = client.Disconnect(context.Background())
	})

	return collection
}

// Тестируем успешное получение детали
func TestRepository_Get_Success(t *testing.T) {
	ctx := context.Background()
	collection := setupTestCollection(t)
	repo := part.NewPartRepository(collection)
	testUUID := "00000000-0000-0000-0000-000000000001"

	_, err := collection.InsertOne(ctx, bson.M{
		"uuid":  testUUID,
		"name":  "Тестовый двигатель",
		"price": 500.0,
	})

	require.NoError(t, err)

	res, err := repo.Get(ctx, testUUID)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, "Тестовый двигатель", res.Name)
	assert.Equal(t, 500.0, res.Price)
}

// Тестируем сценарий, когда деталь отсутствует
func TestRepository_Get_NotFound(t *testing.T) {
	ctx := context.Background()
	collection := setupTestCollection(t)
	repo := part.NewPartRepository(collection)

	_, err := repo.Get(ctx, "unknown-uuid")

	assert.ErrorIs(t, err, model.ErrPartNotFound)
}

// Тестируем получение полного списка без фильтров
func TestRepository_List_All(t *testing.T) {
	ctx := context.Background()
	collection := setupTestCollection(t)
	repo := part.NewPartRepository(collection)

	res, err := repo.List(ctx, nil)

	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}

// Тестируем работу фильтров и вспомогательных функций
func TestRepository_List_WithFilters(t *testing.T) {
	ctx := context.Background()
	collection := setupTestCollection(t)
	repo := part.NewPartRepository(collection)
	testUUID := "00000000-0000-0000-0000-000000000001"

	filterMatch := &model.PartsFilter{
		UUIDs: []string{testUUID},
	}
	resMatch, err := repo.List(ctx, filterMatch)
	assert.NoError(t, err)
	assert.Len(t, resMatch, 1)

	filterName := &model.PartsFilter{
		Names: []string{"Тестовый двигатель"},
	}
	resName, err := repo.List(ctx, filterName)
	assert.NoError(t, err)
	assert.Len(t, resName, 1)

	filterMiss := &model.PartsFilter{
		Names: []string{"Несуществующая деталь для космолета"},
	}
	resMiss, err := repo.List(ctx, filterMiss)
	assert.NoError(t, err)
	assert.Len(t, resMiss, 0)
}
