package model

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusSucceeded PaymentStatus = "SUCCEEDED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
)

type Payment struct {
	PaymentUUID uuid.UUID
	OrderUUID   uuid.UUID
	Amount      float64
	Method      PaymentMethod
	Status      PaymentStatus
	CreatedAt   time.Time
}
