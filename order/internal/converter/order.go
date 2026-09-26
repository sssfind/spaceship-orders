package converter

import (
	"order/internal/model"

	orderV1 "spaceship-orders/shared/pkg/openapi/order/v1"
)

func OrderToDto(order *model.Order) *orderV1.OrderDto {
	if order == nil {
		return nil
	}

	items := make([]orderV1.OrderItemDto, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, orderV1.OrderItemDto{
			PartUUID:  item.PartUUID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		})
	}

	dto := &orderV1.OrderDto{
		OrderUUID:  order.OrderUUID,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartUUIDs,
		Items:      items,
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

func PartToDto(part *model.Part) *orderV1.PartDto {
	if part == nil {
		return nil
	}
	return &orderV1.PartDto{
		PartUUID: part.UUID,
		Name:     part.Name,
		Price:    part.Price,
		Category: part.Category,
		InStock:  part.InStock,
	}
}

func PaymentToDto(payment *model.Payment) *orderV1.PaymentDto {
	if payment == nil {
		return nil
	}
	return &orderV1.PaymentDto{
		PaymentUUID: payment.PaymentUUID,
		OrderUUID:   payment.OrderUUID,
		Amount:      payment.Amount,
		Method:      orderV1.PaymentMethod(payment.Method),
		Status:      orderV1.PaymentStatus(payment.Status),
	}
}
