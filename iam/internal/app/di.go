package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	apiV1 "iam/internal/api/auth/v1"
	"iam/internal/config"
	"iam/internal/repository"
	repoSession "iam/internal/repository/session"
	repoUser "iam/internal/repository/user"
	"iam/internal/service"
	authImpl "iam/internal/service/auth"
	"platform/pkg/closer"

	iamV1 "spaceship-orders/shared/pkg/openapi/iam/v1"
)

type serviceProvider struct {
	cfg *config.Config

	dbPool      *pgxpool.Pool
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	authSrv     service.AuthService
	apiHandler  iamV1.Handler
}

func newServiceProvider(cfg *config.Config) *serviceProvider {
	return &serviceProvider{cfg: cfg}
}

func (sp *serviceProvider) DBPool(ctx context.Context) (*pgxpool.Pool, error) {
	if sp.dbPool == nil {
		pool, err := pgxpool.New(ctx, sp.cfg.Dsn())
		if err != nil {
			return nil, fmt.Errorf("failed to connect to postgres: %w", err)
		}
		sp.dbPool = pool

		closer.AddNamed("postgres_pool", func(_ context.Context) error {
			sp.dbPool.Close()
			return nil
		})
	}
	return sp.dbPool, nil
}

func (sp *serviceProvider) UserRepository(ctx context.Context) (repository.UserRepository, error) {
	if sp.userRepo == nil {
		pool, err := sp.DBPool(ctx)
		if err != nil {
			return nil, err
		}
		sp.userRepo = repoUser.NewRepository(pool)
	}
	return sp.userRepo, nil
}

func (sp *serviceProvider) SessionRepository(ctx context.Context) (repository.SessionRepository, error) {
	if sp.sessionRepo == nil {
		pool, err := sp.DBPool(ctx)
		if err != nil {
			return nil, err
		}
		sp.sessionRepo = repoSession.NewRepository(pool)
	}
	return sp.sessionRepo, nil
}

func (sp *serviceProvider) AuthService(ctx context.Context) (service.AuthService, error) {
	if sp.authSrv == nil {
		users, err := sp.UserRepository(ctx)
		if err != nil {
			return nil, err
		}
		sessions, err := sp.SessionRepository(ctx)
		if err != nil {
			return nil, err
		}
		sp.authSrv = authImpl.NewService(users, sessions, sp.cfg.TTL())
	}
	return sp.authSrv, nil
}

func (sp *serviceProvider) APIHandler(ctx context.Context) (iamV1.Handler, error) {
	if sp.apiHandler == nil {
		srv, err := sp.AuthService(ctx)
		if err != nil {
			return nil, err
		}
		sp.apiHandler = apiV1.NewAPI(srv)
	}
	return sp.apiHandler, nil
}
