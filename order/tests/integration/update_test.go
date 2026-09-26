//go:build integration

package integration

import (
	"github.com/google/uuid"

	"order/internal/model"
)

func (s *OrderTestSuite) TestUpdateStatus_Pay() {
	orderUUID := uuid.New()
	err := s.repo.Create(s.ctx, &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   uuid.New(),
		PartUUIDs:  []uuid.UUID{uuid.New()},
		TotalPrice: 100,
		Status:     model.StatusPendingPayment,
	})
	s.Require().NoError(err)

	txUUID := uuid.New().String()
	err = s.repo.UpdateStatus(s.ctx, orderUUID.String(), model.StatusPaid, txUUID, model.MethodCard)
	s.Require().NoError(err)

	got, err := s.repo.Get(s.ctx, orderUUID.String())
	s.Require().NoError(err)
	s.Equal(model.StatusPaid, got.Status)
	s.Require().NotNil(got.TransactionUUID)
	s.Equal(txUUID, got.TransactionUUID.String())
	s.Require().NotNil(got.PaymentMethod)
	s.Equal(model.MethodCard, *got.PaymentMethod)
}

func (s *OrderTestSuite) TestUpdateStatus_Cancel() {
	orderUUID := uuid.New()
	err := s.repo.Create(s.ctx, &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   uuid.New(),
		PartUUIDs:  []uuid.UUID{uuid.New()},
		TotalPrice: 100,
		Status:     model.StatusPendingPayment,
	})
	s.Require().NoError(err)

	err = s.repo.UpdateStatus(s.ctx, orderUUID.String(), model.StatusCancelled, "", "")
	s.Require().NoError(err)

	got, err := s.repo.Get(s.ctx, orderUUID.String())
	s.Require().NoError(err)
	s.Equal(model.StatusCancelled, got.Status)
}

func (s *OrderTestSuite) TestUpdateStatus_NotFound() {
	err := s.repo.UpdateStatus(s.ctx, uuid.New().String(), model.StatusPaid, uuid.New().String(), model.MethodSbp)
	s.ErrorIs(err, model.ErrOrderNotFound)
}
