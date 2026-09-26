package order

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"order/internal/model"
)

func (s *ServiceSuite) TestDeleteOrder_Success() {
	orderUUID := uuid.New()
	s.repo.EXPECT().
		Delete(mock.Anything, orderUUID.String()).
		Return(nil).
		Once()

	err := s.svc.DeleteOrder(s.ctx, orderUUID)
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestDeleteOrder_NotFound() {
	orderUUID := uuid.New()
	s.repo.EXPECT().
		Delete(mock.Anything, orderUUID.String()).
		Return(model.ErrOrderNotFound).
		Once()

	err := s.svc.DeleteOrder(s.ctx, orderUUID)
	s.ErrorIs(err, model.ErrOrderNotFound)
}
