package data

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/arda-labs/arda/arda-be-go/pkg/database"
	"github.com/arda-labs/arda/arda-be-go/services/mdm-service/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/google/wire"
	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewMdmRepo)

type Data struct {
	db *database.Database
}

func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	if c.Database == nil || c.Database.Source == "" {
		return nil, nil, fmt.Errorf("database source is required")
	}

	db, cleanup, err := database.NewPool(c.Database.Source, logger)
	if err != nil {
		return nil, nil, err
	}

	if err := runMigrations(c.Database.Source); err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("running migrations: %w", err)
	}

	return &Data{db: db}, cleanup, nil
}

func (d *Data) DB(ctx context.Context) *database.Database {
	return d.db
}

func runMigrations(dsn string) error {
	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("creating migration source: %w", err)
	}

	db, err := sql.Open("pgx", dsn)
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
