package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserUUID     uuid.UUID
	Login        string
	PasswordHash string
	CreatedAt    time.Time
}

type Session struct {
	SessionUUID uuid.UUID
	UserUUID    uuid.UUID
	ExpiresAt   time.Time
	CreatedAt   time.Time
}
