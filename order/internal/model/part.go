package model

import (
	"time"

	"github.com/google/uuid"
)

type Part struct {
	UUID      uuid.UUID
	Name      string
	Price     float64
	Category  string
	InStock   int
	CreatedAt time.Time
}
