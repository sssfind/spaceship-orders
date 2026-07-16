package env

import (
	"errors"
	"fmt"
	"os"
)

type postgres struct {
	dsn          string
	migrationDir string
}

func NewPostgresConfig() (*postgres, error) {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	if host == "" || port == "" || user == "" || password == "" || dbName == "" {
		return nil, errors.New("critical postgres configuration variables are missing")
	}

	sslmode := os.Getenv("POSTGRES_SSL_MODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	migrationDir := os.Getenv("MIGRATION_DIRECTORY")
	if migrationDir == "" {
		migrationDir = "migrations"
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbName, sslmode)

	return &postgres{
		dsn:          dsn,
		migrationDir: migrationDir,
	}, nil
}

func (cfg *postgres) Dsn() string {
	return cfg.dsn
}

func (cfg *postgres) MigrationDir() string {
	return cfg.migrationDir
}
