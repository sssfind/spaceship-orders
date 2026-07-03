package converter

import (
	"order/internal/model"
	orderV1 "spaceship-orders/shared/pkg/openapi/order/v1"
)

func OrderToDto(order *model.Order) *orderV1.OrderDto {
	if order == nil {
		return nil
	}

	dto := &orderV1.OrderDto{
		OrderUUID:  order.OrderUUID,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartUUIDs,
		TotalPrice: order.TotalPrice,
		Status:     orderV1.OrderStatus(order.Status),
	}

	if order.TransactionUUID != nil {
		dto.TransactionUUID = orderV1.NewOptUUID(*order.TransactionUUID)
	}
	if order.PaymentMethod != nil {
		dto.PaymentMethod = orderV1.NewOptNilPaymentMethod(orderV1.PaymentMethod(*order.PaymentMethod))
	}

	return dto
}
