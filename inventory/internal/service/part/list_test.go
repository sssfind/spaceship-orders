package part_test

import (
	"inventory/internal/model"

	"github.com/stretchr/testify/mock"
)

func (s *PartServiceTestSuite) TestListParts_Success() {
	mockParts := []model.Part{
		{UUID: "1", Name: "Деталь 1", Price: 10},
		{UUID: "2", Name: "Деталь 2", Price: 20},
	}

	s.repoMock.On("List", mock.Anything, mock.AnythingOfType("*model.PartsFilter")).Return(mockParts, nil).Once()

	res, err := s.service.ListParts(s.ctx, &model.PartsFilter{})

	s.NoError(err)
	s.Len(res, 2)
	s.Equal("Деталь 1", res[0].Name)
}
