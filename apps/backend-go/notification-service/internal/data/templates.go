package data

import (
	"context"
	"errors"

	"github.com/arda-labs/arda/arda-be-go/pkg/pagination"
	"github.com/arda-labs/arda/arda-be-go/services/notification-service/internal/biz"
	"github.com/jackc/pgx/v5"
)

func (r *NotificationRepo) ListTemplates(ctx context.Context, filter biz.PageFilter) ([]*biz.NotificationTemplate, string, error) {
	page := pagination.Normalize(filter.PageSize, filter.PageToken)
	rows, err := r.data.db.Pool.Query(ctx, `
		SELECT id::text, code, name, category, default_channel, description, status, created_at, updated_at
		FROM notification_templates
		WHERE deleted_at IS NULL
		  AND ($1 = '' OR status = $1)
		  AND ($2 = '' OR code ILIKE '%' || $2 || '%' OR name ILIKE '%' || $2 || '%' OR category ILIKE '%' || $2 || '%')
		ORDER BY category ASC, code ASC
		LIMIT $3 OFFSET $4`, filter.Status, filter.Keyword, page.Limit+1, page.Offset)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()
	list, err := scanTemplates(rows)
	if err != nil {
		return nil, "", err
	}
	next := pagination.NextOffsetToken(len(list), page.Limit, page.Offset)
	if len(list) > page.Limit {
		list = list[:page.Limit]
	}
	return list, next, nil
}

func (r *NotificationRepo) GetTemplate(ctx context.Context, id string) (*biz.NotificationTemplate, error) {
	row := r.data.db.Pool.QueryRow(ctx, `
		SELECT id::text, code, name, category, default_channel, description, status, created_at, updated_at
		FROM notification_templates
		WHERE id = $1 AND deleted_at IS NULL`, id)
	return scanTemplate(row)
}

func (r *NotificationRepo) CreateTemplate(ctx context.Context, item *biz.NotificationTemplate) (*biz.NotificationTemplate, error) {
	err := r.data.db.Pool.QueryRow(ctx, `
		INSERT INTO notification_templates (code, name, category, default_channel, description, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id::text`,
		item.Code, item.Name, item.Category, item.DefaultChannel, item.Description, item.Status).Scan(&item.ID)
	if err != nil {
		return nil, err
	}
	return r.GetTemplate(ctx, item.ID)
}

func (r *NotificationRepo) UpdateTemplate(ctx context.Context, item *biz.NotificationTemplate) (*biz.NotificationTemplate, error) {
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE notification_templates
		SET code=$2, name=$3, category=$4, default_channel=$5, description=$6, status=$7, updated_at=now()
		WHERE id=$1 AND deleted_at IS NULL`,
		item.ID, item.Code, item.Name, item.Category, item.DefaultChannel, item.Description, item.Status)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.GetTemplate(ctx, item.ID)
}

func (r *NotificationRepo) DeleteTemplate(ctx context.Context, id string) error {
	tag, err := r.data.db.Pool.Exec(ctx, `UPDATE notification_templates SET status='DELETED', deleted_at=now(), updated_at=now() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return biz.ErrNotFound
	}
	return nil
}

func scanTemplates(rows pgx.Rows) ([]*biz.NotificationTemplate, error) {
	var list []*biz.NotificationTemplate
	for rows.Next() {
		item, err := scanTemplate(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func scanTemplate(row pgx.Row) (*biz.NotificationTemplate, error) {
	item := &biz.NotificationTemplate{}
	err := row.Scan(&item.ID, &item.Code, &item.Name, &item.Category, &item.DefaultChannel, &item.Description, &item.Status, &item.CreatedAt, &item.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
}
