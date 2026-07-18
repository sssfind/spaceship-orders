package model

import "errors"

var (
	// ErrUserNotFound пользователь не найден в системе
	ErrUserNotFound = errors.New("user not found")

	// ErrUserAlreadyExists логин или email уже заняты
	ErrUserAlreadyExists = errors.New("user with this login or email already exists")

	// ErrInvalidCredentials неверный логин или пароль при входе
	ErrInvalidCredentials = errors.New("invalid login or password")

	// ErrSessionNotFound сессия отсутствует в Redis или ее TTL истек
	ErrSessionNotFound = errors.New("session not found or expired")
)
