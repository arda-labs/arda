package data

import (
	"context"
	"errors"

	"github.com/arda-labs/arda/arda-be-go/services/notification-service/internal/biz"
	"github.com/jackc/pgx/v5"
)

func (r *NotificationRepo) ListProviderConfigs(ctx context.Context, filter biz.ProviderConfigFilter) ([]*biz.ProviderConfig, error) {
	rows, err := r.data.db.Pool.Query(ctx, providerConfigSelect()+`
		WHERE ($1 = '' OR channel = $1)
		  AND ($2 = '' OR status = $2)
		ORDER BY channel ASC, priority ASC, code ASC`, filter.Channel, filter.Status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProviderConfigs(rows)
}

func (r *NotificationRepo) UpsertProviderConfig(ctx context.Context, item *biz.ProviderConfig) (*biz.ProviderConfig, error) {
	err := r.data.db.Pool.QueryRow(ctx, `
		INSERT INTO notification_provider_configs (code, channel, name, priority, rate_limit_per_minute, options_json, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (code) DO UPDATE
		SET channel=EXCLUDED.channel,
		    name=EXCLUDED.name,
		    priority=EXCLUDED.priority,
		    rate_limit_per_minute=EXCLUDED.rate_limit_per_minute,
		    options_json=EXCLUDED.options_json,
		    status=EXCLUDED.status,
		    updated_at=now()
		RETURNING id::text`,
		item.Code, item.Channel, item.Name, item.Priority, item.RateLimitPerMinute, item.OptionsJSON, item.Status).Scan(&item.ID)
	if err != nil {
		return nil, err
	}
	return r.getProviderConfig(ctx, item.ID)
}

func (r *NotificationRepo) getProviderConfig(ctx context.Context, id string) (*biz.ProviderConfig, error) {
	row := r.data.db.Pool.QueryRow(ctx, providerConfigSelect()+` WHERE id=$1`, id)
	return scanProviderConfig(row)
}

func providerConfigSelect() string {
	return `SELECT id::text, code, channel, name, priority, rate_limit_per_minute, options_json, status, created_at, updated_at FROM notification_provider_configs`
}

func scanProviderConfigs(rows pgx.Rows) ([]*biz.ProviderConfig, error) {
	var list []*biz.ProviderConfig
	for rows.Next() {
		item, err := scanProviderConfig(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func scanProviderConfig(row pgx.Row) (*biz.ProviderConfig, error) {
	item := &biz.ProviderConfig{}
	err := row.Scan(&item.ID, &item.Code, &item.Channel, &item.Name, &item.Priority, &item.RateLimitPerMinute, &item.OptionsJSON, &item.Status, &item.CreatedAt, &item.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
}
