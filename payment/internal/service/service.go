package service

import (
	"context"

	"payment/internal/model"
)

type PaymentService interface {
	ProcessPayment(ctx context.Context, orderUUID, userUUID string, method model.PaymentMethod) (string, error)
}
