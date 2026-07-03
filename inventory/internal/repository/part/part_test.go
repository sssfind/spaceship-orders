package part_test

import (
	"context"
	"testing"

	"inventory/internal/model"
	"inventory/internal/repository/part"

	"github.com/stretchr/testify/assert"
)

// Тестируем успешное получение детали
func TestRepository_Get_Success(t *testing.T) {
	ctx := context.Background()
	repo := part.NewPartRepository()
	testUUID := "00000000-0000-0000-0000-000000000001"

	res, err := repo.Get(ctx, testUUID)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Тестовый двигатель", res.Name)
	assert.Equal(t, 500.0, res.Price)
}

// Тестируем сценарий, когда деталь отсутствует
func TestRepository_Get_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := part.NewPartRepository()

	_, err := repo.Get(ctx, "unknown-uuid")

	assert.ErrorIs(t, err, model.ErrPartNotFound)
}

// Тестируем получение полного списка без фильтров
func TestRepository_List_All(t *testing.T) {
	ctx := context.Background()
	repo := part.NewPartRepository()

	res, err := repo.List(ctx, nil)

	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}

// Тестируем работу фильтров и вспомогательных функций
func TestRepository_List_WithFilters(t *testing.T) {
	ctx := context.Background()
	repo := part.NewPartRepository()
	testUUID := "00000000-0000-0000-0000-000000000001"

	// фильтр по UUID совпадает
	filterMatch := &model.PartsFilter{
		UUIDs: []string{testUUID},
	}
	resMatch, err := repo.List(ctx, filterMatch)
	assert.NoError(t, err)
	assert.Len(t, resMatch, 1)

	// фильтр по имени совпадает
	filterName := &model.PartsFilter{
		Names: []string{"Тестовый двигатель"},
	}
	resName, err := repo.List(ctx, filterName)
	assert.NoError(t, err)
	assert.Len(t, resName, 1)

	// Фильтр по имени не совпадает
	filterMiss := &model.PartsFilter{
		Names: []string{"Несуществующая деталь для космолета"},
	}
	resMiss, err := repo.List(ctx, filterMiss)
	assert.NoError(t, err)
	assert.Len(t, resMiss, 0)

	// фильтр по категории совпадает
	filterCategory := &model.PartsFilter{
		Categories: []model.Category{model.Category(1)}, // CategoryEngine
	}
	resCat, err := repo.List(ctx, filterCategory)
	assert.NoError(t, err)
	assert.Len(t, resCat, 1)
}
