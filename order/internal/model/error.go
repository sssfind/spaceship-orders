package model

import "errors"

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrOrderAlreadyPaid   = errors.New("order is already paid")
	ErrInvalidOrderStatus = errors.New("invalid order status")
	ErrPartNotFound       = errors.New("part not found")
	ErrPartAlreadyExists  = errors.New("part already exists")
	ErrInsufficientStock  = errors.New("insufficient stock")
	ErrPaymentNotFound    = errors.New("payment not found")
	ErrEmptyPartsList     = errors.New("parts list is empty")
)
