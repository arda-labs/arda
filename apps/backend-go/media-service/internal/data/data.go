package data

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type Data struct {
	db *sql.DB
}

func NewData(databaseURL string) (*Data, func(), error) {
	if databaseURL == "" {
		return nil, nil, fmt.Errorf("database url is required")
	}
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, nil, err
	}
	if err := runMigrations(databaseURL); err != nil {
		db.Close()
		return nil, nil, err
	}
	return &Data{db: db}, func() { _ = db.Close() }, nil
}

func (d *Data) DB() *sql.DB {
	return d.db
}

func (d *Data) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

func runMigrations(databaseURL string) error {
	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("creating migration source: %w", err)
	}
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return err
	}
	defer db.Close()
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
