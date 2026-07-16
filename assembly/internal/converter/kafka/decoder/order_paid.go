package decoder

import (
	"assembly/internal/model"
	pbEvents "spaceship-orders/shared/pkg/proto/events/v1"
)

func ToOrderPaidModel(protoEvent *pbEvents.OrderPaidEvent) model.OrderPaid {
	return model.OrderPaid{
		EventUUID:       protoEvent.GetEventUuid(),
		OrderUUID:       protoEvent.GetOrderUuid(),
		UserUUID:        protoEvent.GetUserUuid(),
		PaymentMethod:   protoEvent.GetPaymentMethod(),
		TransactionUUID: protoEvent.GetTransactionUuid(),
	}
}
