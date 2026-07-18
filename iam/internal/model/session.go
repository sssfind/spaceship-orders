package model

import "time"

// Session описывает активную сессию пользователя
type Session struct {
	SessionUUID string
	UserUUID    string
	CreatedAt   time.Time
}
