package model

import "time"

type UserRepoModel struct {
	ID                  int64     `db:"id"`
	UUID                string    `db:"uuid"`
	Login               string    `db:"login"`
	PasswordHash        string    `db:"password_hash"`
	Email               string    `db:"email"`
	NotificationMethods []byte    `db:"notification_methods"` // Храним как JSONB в Postgres
	CreatedAt           time.Time `db:"created_at"`
}
