package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	healthAPI "order/internal/api/health"
	iamclient "order/internal/client/iam"
	"order/internal/config"
	customMiddleware "order/internal/middleware"
	"platform/pkg/closer"
	"platform/pkg/logger"
	platformMigrator "platform/pkg/migrator/pg"
	"platform/pkg/tracing"

	orderV1 "spaceship-orders/shared/pkg/openapi/order/v1"
)

type App struct {
	serviceProvider *serviceProvider
	httpServer      *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	// Загружаем конфигурацию из переменных окружения
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Инициализируем платформенный логгер с поддержкой мульти-кора
	err = logger.InitWithConfig(logger.Config{
		Level:                 cfg.LogLevel(),
		AsJSON:                cfg.LogAsJSON(),
		ServiceName:           cfg.ServiceName(),
		Outputs:               cfg.Outputs(),
		OtelCollectorEndpoint: cfg.OtelCollectorEndpoint(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init logger: %w", err)
	}

	closer.AddNamed("logger", func(c context.Context) error {
		return logger.Shutdown(c)
	})

	err = tracing.InitTracer(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init tracer: %w", err)
	}

	closer.AddNamed("tracer", func(c context.Context) error {
		return tracing.ShutdownTracer(c)
	})

	a.serviceProvider = newServiceProvider(cfg)

	err = a.initDependencies(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initDependencies(ctx context.Context) error {
	dbMigrator := platformMigrator.NewMigrator(
		a.serviceProvider.cfg.Dsn(),
		a.serviceProvider.cfg.MigrationDir(),
	)
	if err := dbMigrator.Up(); err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}
	logger.Info(ctx, "Database migrations successfully applied")

	// Собираем хендлеры
	apiHandler, err := a.serviceProvider.APIHandler(ctx)
	if err != nil {
		return err
	}

	orderServer, err := orderV1.NewServer(apiHandler)
	if err != nil {
		return fmt.Errorf("failed to create openapi server: %w", err)
	}

	healthHandler := healthAPI.NewHandler()

	iamBaseURL := a.serviceProvider.cfg.IAMHTTPConfig.BaseURL()
	iamTarget, err := url.Parse(iamBaseURL)
	if err != nil {
		return fmt.Errorf("invalid IAM_HTTP_BASE_URL: %w", err)
	}
	iamProxy := httputil.NewSingleHostReverseProxy(iamTarget)
	iamClient := iamclient.NewClient(iamBaseURL)
	authMW := customMiddleware.NewAuthMiddleware(iamClient)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(customMiddleware.RequestLogger)

	r.Get("/health", healthHandler.HealthCheck)
	r.Handle("/metrics", promhttp.Handler())

	if frontendDir := resolveFrontendDir(); frontendDir != "" {
		mountFrontend(r, frontendDir)
		logger.Info(ctx, fmt.Sprintf("frontend mounted from %s", frontendDir))
	} else {
		logger.Info(ctx, "frontend directory not found, UI is disabled")
	}

	// Same-origin auth for the UI: Order proxies /api/v1/auth/* to IAM.
	r.Handle("/api/v1/auth/*", iamProxy)

	r.Group(func(pr chi.Router) {
		pr.Use(authMW.Handle)
		pr.Mount("/", orderServer)
	})

	a.httpServer = &http.Server{
		Addr:              a.serviceProvider.cfg.OrderHttpConfig.Address(),
		Handler:           r,
		ReadHeaderTimeout: a.serviceProvider.cfg.ReadTimeout(),
	}

	closer.AddNamed("http_server", func(c context.Context) error {
		return a.httpServer.Shutdown(c)
	})

	return nil
}

func resolveFrontendDir() string {
	candidates := []string{"frontend", filepath.Join("..", "frontend"), "/app/frontend"}
	if dir := os.Getenv("FRONTEND_DIR"); dir != "" {
		candidates = append([]string{dir}, candidates...)
	}

	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			abs, err := filepath.Abs(candidate)
			if err != nil {
				return candidate
			}
			return abs
		}
	}

	return ""
}

func mountFrontend(r chi.Router, dir string) {
	serve := func(name string) http.HandlerFunc {
		path := filepath.Join(dir, name)
		return func(w http.ResponseWriter, req *http.Request) {
			http.ServeFile(w, req, path)
		}
	}

	r.Get("/", serve("index.html"))
	r.Get("/styles.css", serve("styles.css"))
	r.Get("/app.js", serve("app.js"))
}

func (a *App) Run() error {
	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	closer.SetLogger(logger.Logger())

	ctx := context.Background()

	logger.Info(ctx, fmt.Sprintf("order HTTP Server успешно запущен на %s", a.httpServer.Addr))

	if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server runtime error: %w", err)
	}

	return nil
}
