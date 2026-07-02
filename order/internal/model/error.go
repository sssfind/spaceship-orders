package model

import "errors"

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrOrderAlreadyPaid   = errors.New("order is already paid")
	ErrInvalidOrderStatus = errors.New("invalid order status")
)
