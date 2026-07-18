package session

import (
	"context"
	"encoding/json"
	"fmt"
	"iam/internal/model"
	"iam/internal/repository/converter"
	"time"
)

func (r *sessionRepo) Create(ctx context.Context, s *model.Session, ttl time.Duration) error {
	repoModel := converter.ToRepoFromSession(s)

	data, err := json.Marshal(repoModel)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("session:%s", s.SessionUUID)

	// Используем специализированный метод платформенного клиента для TTL
	return r.cache.SetWithTTL(ctx, key, data, ttl)
}
