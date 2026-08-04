package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
	"iam/internal/model"
	"iam/internal/repository/converter"
	repoModel "iam/internal/repository/model"
)

func (r *sessionRepo) Get(ctx context.Context, sessionUUID string) (*model.Session, error) {
	key := fmt.Sprintf("session:%s", sessionUUID)

	// Платформенный клиент возвращает слайс байт напрямую
	data, err := r.cache.Get(ctx, key)
	if err != nil {
		// Если ключ в Redis отсутствует, библиотека redigo возвращает ErrNil
		if errors.Is(err, redigo.ErrNil) {
			return nil, model.ErrSessionNotFound
		}
		return nil, err
	}

	var rm repoModel.SessionRepoModel
	if err := json.Unmarshal(data, &rm); err != nil {
		return nil, err
	}

	return converter.ToSessionFromRepo(&rm), nil
}
