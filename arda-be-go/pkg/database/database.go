package database

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func NewPool(source string, logger log.Logger) (*Database, func(), error) {
	l := log.NewHelper(logger)

	config, err := pgxpool.ParseConfig(source)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Default pool settings
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		return nil, nil, fmt.Errorf("failed to ping database: %w", err)
	}

	l.Info("database connection pool initialized")

	cleanup := func() {
		pool.Close()
		l.Info("database connection pool closed")
	}

	return &Database{Pool: pool}, cleanup, nil
}

// ExecInTransaction runs a function within a transaction and sets the tenant ID for RLS
func (db *Database) ExecInTransaction(ctx context.Context, tenantID string, fn func(ctx context.Context, tx pgx.Tx) error) error {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if tenantID != "" {
		if _, err := tx.Exec(ctx, "SET LOCAL app.current_tenant_id = $1", tenantID); err != nil {
			return err
		}
	}

	if err := fn(ctx, tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
