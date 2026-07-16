package part_test

import (
	"github.com/stretchr/testify/mock"
	"inventory/internal/model"
)

func (s *PartServiceTestSuite) TestGetPart_Success() {
	partUUID := "00000000-0000-0000-0000-000000000001"
	mockPart := &model.Part{
		UUID:  partUUID,
		Name:  "Тестовый двигатель",
		Price: 500.0,
	}

	s.repoMock.On("Get", mock.Anything, partUUID).Return(mockPart, nil).Once()

	res, err := s.service.GetPart(s.ctx, partUUID)

	s.NoError(err)
	s.NotNil(res)
	s.Equal(partUUID, res.UUID)
}

func (s *PartServiceTestSuite) TestGetPart_NotFound() {
	partUUID := "non-existent-uuid"

	s.repoMock.On("Get", mock.Anything, partUUID).Return(nil, model.ErrPartNotFound).Once()

	_, err := s.service.GetPart(s.ctx, partUUID)

	s.ErrorIs(err, model.ErrPartNotFound)
}
