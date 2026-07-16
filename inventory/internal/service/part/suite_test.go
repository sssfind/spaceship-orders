package part_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	repoMocks "inventory/internal/repository/mocks"
	"inventory/internal/service"
	partService "inventory/internal/service/part"
)

type PartServiceTestSuite struct {
	suite.Suite
	ctx      context.Context
	repoMock *repoMocks.PartRepository
	service  service.PartService
}

func (s *PartServiceTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.repoMock = repoMocks.NewPartRepository(s.T())
	s.service = partService.NewService(s.repoMock)
}

func TestPartServiceSuite(t *testing.T) {
	suite.Run(t, new(PartServiceTestSuite))
}
