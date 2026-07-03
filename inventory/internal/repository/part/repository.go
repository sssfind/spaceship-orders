package part

import (
	"sync"

	"inventory/internal/repository"
	repoModel "inventory/internal/repository/model"
)

type repo struct {
	mu    sync.RWMutex
	parts map[string]*repoModel.Part
}

func NewPartRepository() repository.PartRepository {
	r := &repo{
		parts: make(map[string]*repoModel.Part),
	}
	r.initTestData() // Заполняем тестовыми данными из init.go
	return r
}
