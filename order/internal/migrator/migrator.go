package migrator

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type Migrator struct {
	dsn           string
	migrationsDir string
}

func NewMigrator(dsn, migrationsDir string) *Migrator {
	return &Migrator{
		dsn:           dsn,
		migrationsDir: migrationsDir,
	}
}

func (m *Migrator) Up() error {
	db, err := sql.Open("pgx", m.dsn)
	if err != nil {
		return fmt.Errorf("opening database connection: %w", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("setting postgres dialect: %w", err)
	}

	if err := goose.Up(db, m.migrationsDir); err != nil {
		return fmt.Errorf("running migrations: %w", err)
	}

	return nil
}
