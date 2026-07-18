package session

import (
	"context"
	"iam/internal/repository"
	"time"
)

type RedisClient interface {
	SetWithTTL(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
}

type sessionRepo struct {
	cache RedisClient
}

func NewSessionRepository(cache RedisClient) repository.SessionRepository {
	return &sessionRepo{
		cache: cache,
	}
}
