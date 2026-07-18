package env

import (
	"errors"
	"fmt"
	"os"
)

type postgresConfig struct {
	dsn          string
	migrationDir string
}

func NewPostgresConfig() (*postgresConfig, error) {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		return nil, errors.New("POSTGRES_HOST is not set")
	}
	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		return nil, errors.New("POSTGRES_PORT is not set")
	}
	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		return nil, errors.New("POSTGRES_USER is not set")
	}
	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		return nil, errors.New("POSTGRES_PASSWORD is not set")
	}
	dbname := os.Getenv("POSTGRES_DB")
	if dbname == "" {
		return nil, errors.New("POSTGRES_DB is not set")
	}
	sslMode := os.Getenv("POSTGRES_SSL_MODE")
	if sslMode == "" {
		sslMode = "disable"
	}
	migrationDir := os.Getenv("MIGRATION_DIRECTORY")
	if migrationDir == "" {
		migrationDir = "migrations"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslMode)

	return &postgresConfig{
		dsn:          dsn,
		migrationDir: migrationDir,
	}, nil
}

func (p *postgresConfig) DSN() string {
	return p.dsn
}

func (p *postgresConfig) MigrationDir() string {
	return p.migrationDir
}
