package data

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/arda-labs/arda/arda-be-go/pkg/database"
	"github.com/arda-labs/arda/arda-be-go/pkg/redis"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"github.com/arda-labs/arda/arda-be-go/services/iam-service/internal/conf"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

var ProviderSet = wire.NewSet(
	NewData,
	NewUserRepo,
	NewTenantRepo,
	NewMembershipRepo,
	NewRoleRepo,
	NewPermissionRepo,
	NewAuditRepo,
	NewPermissionCache,
	NewMenuRepo,
)

type Data struct {
	db  *database.Database
	rdb *redis.RedisClient
}

func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	db, dbCleanup, err := database.NewPool(c.Database.Source, logger)
	if err != nil {
		return nil, nil, err
	}

	if err := runMigrations(c.Database.Source); err != nil {
		log.NewHelper(logger).Warnf("migration warning: %v", err)
	}

	rdb, redisCleanup, err := redis.NewRedis(
		c.Redis.Addr,
		int(c.Redis.GetDb()),
		"", // password if needed
		logger,
	)
	if err != nil {
		dbCleanup()
		return nil, nil, err
	}

	cleanup := func() {
		dbCleanup()
		redisCleanup()
	}

	return &Data{db: db, rdb: rdb}, cleanup, nil
}

func (d *Data) DB(ctx context.Context) *database.Database {
	return d.db
}

func runMigrations(dsn string) error {
	d, err := iofs.New(migrationsFS, "migrations")
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
	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
