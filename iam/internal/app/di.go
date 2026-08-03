package app

import (
	"context"
	"database/sql"
	"fmt"

	// Драйверы и внешние либы
	redigo "github.com/gomodule/redigo/redis"
	_ "github.com/lib/pq" // Обязательно для работы database/sql с Postgres
	"go.uber.org/zap"
	authApi "iam/internal/api/auth/v1"
	userApi "iam/internal/api/user/v1"
	"iam/internal/config"
	"iam/internal/repository"
	sessionRepo "iam/internal/repository/session"
	userRepo "iam/internal/repository/user"
	"iam/internal/service"
	authService "iam/internal/service/auth"
	userService "iam/internal/service/user"
	platformRedis "platform/pkg/cache/redis"
)

type redisLoggerWrapper struct {
	log *zap.Logger
}

func (w *redisLoggerWrapper) Info(ctx context.Context, msg string, fields ...zap.Field) {
	w.log.Info(msg, fields...)
}

func (w *redisLoggerWrapper) Error(ctx context.Context, msg string, fields ...zap.Field) {
	w.log.Error(msg, fields...)
}

// Container инкапсулирует в себе все зависимости приложения
type Container struct {
	config config.Config

	// Инфраструктурные зависимости
	db          *sql.DB
	redisPool   *redigo.Pool
	redisClient sessionRepo.RedisClient // Использован локальный интерфейс репозитория вместо приватного типа платформы

	// Репозитории
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository

	// Сервисы
	userServ service.UserService
	authServ service.AuthService

	// API Хендлеры
	authImpl *authApi.Implementation
	userImpl *userApi.Implementation
}

// NewContainer создает новый экземпляр DI-контейнера
func NewContainer(cfg *config.Config) *Container {
	return &Container{
		config: *cfg,
	}
}

// DB инициализирует пул подключений к PostgreSQL
func (c *Container) DB(ctx context.Context) (*sql.DB, error) {
	if c.db == nil {
		db, err := sql.Open("postgres", c.config.Postgres().DSN())
		if err != nil {
			return nil, fmt.Errorf("failed to open postgres: %w", err)
		}

		// Проверяем жива ли база
		if err := db.PingContext(ctx); err != nil {
			return nil, fmt.Errorf("failed to ping postgres: %w", err)
		}

		c.db = db
	}
	return c.db, nil
}

// RedisPool инициализирует пул соединений redigo
func (c *Container) RedisPool() *redigo.Pool {
	if c.redisPool == nil {
		c.redisPool = &redigo.Pool{
			MaxIdle:     c.config.Redis().MaxIdle(),
			IdleTimeout: c.config.Redis().IdleTimeout(),
			Dial: func() (redigo.Conn, error) {
				return redigo.Dial("tcp", c.config.Redis().Address())
			},
		}
	}
	return c.redisPool
}

// RedisClient собирает платформенный клиент поверх пула redigo
func (c *Container) RedisClient() sessionRepo.RedisClient {
	if c.redisClient == nil {
		zapLogger, _ := zap.NewProduction()

		platformLogger := &redisLoggerWrapper{log: zapLogger}

		// Собираем клиент из платформы
		c.redisClient = platformRedis.NewClient(
			c.RedisPool(),
			platformLogger,
			c.config.Redis().ConnectionTimeout(),
		)
	}
	return c.redisClient
}

// UserRepository синглтон репозитория пользователей
func (c *Container) UserRepository(ctx context.Context) (repository.UserRepository, error) {
	if c.userRepository == nil {
		db, err := c.DB(ctx)
		if err != nil {
			return nil, err
		}
		c.userRepository = userRepo.NewUserRepository(db)
	}
	return c.userRepository, nil
}

// SessionRepository синглтон репозитория сессий
func (c *Container) SessionRepository() repository.SessionRepository {
	if c.sessionRepository == nil {
		c.sessionRepository = sessionRepo.NewSessionRepository(c.RedisClient())
	}
	return c.sessionRepository
}

// UserService синглтон бизнес-логики пользователей
func (c *Container) UserService(ctx context.Context) (service.UserService, error) {
	if c.userServ == nil {
		uRepo, err := c.UserRepository(ctx)
		if err != nil {
			return nil, err
		}
		c.userServ = userService.NewUserService(uRepo)
	}
	return c.userServ, nil
}

// AuthService синглтон бизнес-логики аутентификации
func (c *Container) AuthService(ctx context.Context) (service.AuthService, error) {
	if c.authServ == nil {
		uRepo, err := c.UserRepository(ctx)
		if err != nil {
			return nil, err
		}
		c.authServ = authService.NewAuthService(uRepo, c.SessionRepository(), c.config.Session())
	}
	return c.authServ, nil
}

// AuthImpl возвращает готовый gRPC хендлер для AuthService
func (c *Container) AuthImpl(ctx context.Context) (*authApi.Implementation, error) {
	if c.authImpl == nil {
		aServ, err := c.AuthService(ctx)
		if err != nil {
			return nil, err
		}
		c.authImpl = authApi.NewImplementation(aServ)
	}
	return c.authImpl, nil
}

// UserImpl returns готовый gRPC хендлер для UserService
func (c *Container) UserImpl(ctx context.Context) (*userApi.Implementation, error) {
	if c.userImpl == nil {
		uServ, err := c.UserService(ctx)
		if err != nil {
			return nil, err
		}
		c.userImpl = userApi.NewImplementation(uServ)
	}
	return c.userImpl, nil
}

// Close корректно закрывает все инфраструктурные коннекты при выходе
func (c *Container) Close() error {
	if c.db != nil {
		if err := c.db.Close(); err != nil {
			return fmt.Errorf("error closing postgres pool: %w", err)
		}
	}
	if c.redisPool != nil {
		if err := c.redisPool.Close(); err != nil {
			return fmt.Errorf("error closing redis pool: %w", err)
		}
	}
	return nil
}
