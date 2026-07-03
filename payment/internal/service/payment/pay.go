package payment

import (
	"context"
	"log"

	"github.com/google/uuid"
	"payment/internal/model"
)

func (s *srv) ProcessPayment(ctx context.Context, orderUUID, userUUID string, method model.PaymentMethod) (string, error) {
	txUUID := uuid.NewString()

	log.Printf("Проведена оплата для заказа %s, метод: %s, TX: %s", orderUUID, method, txUUID)

	return txUUID, nil
}
