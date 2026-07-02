package model

import "github.com/google/uuid"

type OrderStatus string
type PaymentMethod string

const (
	StatusPendingPayment OrderStatus = "1"
	StatusPaid           OrderStatus = "PAID"
	StatusCancelled      OrderStatus = "CANCELLED"
)

type Order struct {
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	PartUUIDs       []uuid.UUID
	TotalPrice      float64
	Status          OrderStatus
	TransactionUUID *uuid.UUID
	PaymentMethod   *PaymentMethod
}
