package model

import "github.com/google/uuid"

type OrderItem struct {
	OrderUUID uuid.UUID
	PartUUID  uuid.UUID
	Quantity  int
	UnitPrice float64
}
