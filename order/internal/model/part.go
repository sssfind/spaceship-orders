package model

import "github.com/google/uuid"

type Part struct {
	UUID  uuid.UUID
	Price float64
	Name  string
}
