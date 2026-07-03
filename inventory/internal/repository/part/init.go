package part

import repoModel "inventory/internal/repository/model"

func (r *repo) initTestData() {
	testPartUUID := "00000000-0000-0000-0000-000000000001"
	r.parts[testPartUUID] = &repoModel.Part{
		UUID:     testPartUUID,
		Name:     "Тестовый двигатель",
		Price:    500.0,
		Category: repoModel.Category(1),
	}
}