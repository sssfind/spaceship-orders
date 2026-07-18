package model

type SessionRepoModel struct {
	SessionUUID string `json:"session_uuid"`
	UserUUID    string `json:"user_uuid"`
	CreatedAt   int64  `json:"created_at"` // Unix timestamp для хранения в Redis
}
