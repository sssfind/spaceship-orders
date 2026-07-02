package order

import (
	"context"
	"order/internal/model"
)

type PaymentClient interface {
	Pay(ctx context.Context, orderUUID, userUUID string, method model.PaymentMethod) (string, error)
}
