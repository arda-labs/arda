package data

import (
	"context"
	"errors"

	"github.com/arda-labs/arda/arda-be-go/services/notification-service/internal/biz"
	"github.com/jackc/pgx/v5"
)

func (r *NotificationRepo) ListTemplateVersions(ctx context.Context, filter biz.TemplateVersionFilter) ([]*biz.NotificationTemplateVersion, error) {
	rows, err := r.data.db.Pool.Query(ctx, templateVersionSelect()+`
		WHERE deleted_at IS NULL
		  AND template_id = $1
		  AND ($2 = '' OR channel = $2)
		  AND ($3 = '' OR language = $3)
		  AND ($4 = '' OR status = $4)
		ORDER BY version DESC, channel ASC, language ASC`, filter.TemplateID, filter.Channel, filter.Language, filter.Status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTemplateVersions(rows)
}

func (r *NotificationRepo) CreateTemplateVersion(ctx context.Context, item *biz.NotificationTemplateVersion) (*biz.NotificationTemplateVersion, error) {
	err := r.data.db.Pool.QueryRow(ctx, `
		INSERT INTO notification_template_versions (
			template_id, version, channel, language, subject, body, payload_schema_json,
			approval_status, approved_by, change_note, status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id::text`,
		item.TemplateID, item.Version, item.Channel, item.Language, item.Subject, item.Body, item.PayloadSchemaJSON,
		item.ApprovalStatus, item.ApprovedBy, item.ChangeNote, item.Status).Scan(&item.ID)
	if err != nil {
		return nil, err
	}
	return r.getTemplateVersion(ctx, item.ID)
}

func (r *NotificationRepo) ApproveTemplateVersion(ctx context.Context, id, actor, note string) (*biz.NotificationTemplateVersion, error) {
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE notification_template_versions
		SET approval_status='APPROVED', approved_by=$2, approved_at=now(), change_note=$3, status='ACTIVE', updated_at=now()
		WHERE id=$1 AND deleted_at IS NULL`, id, actor, note)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.getTemplateVersion(ctx, id)
}

func (r *NotificationRepo) getTemplateVersion(ctx context.Context, id string) (*biz.NotificationTemplateVersion, error) {
	row := r.data.db.Pool.QueryRow(ctx, templateVersionSelect()+` WHERE id=$1 AND deleted_at IS NULL`, id)
	return scanTemplateVersion(row)
}

func templateVersionSelect() string {
	return `SELECT id::text, template_id::text, version, channel, language, subject, body, payload_schema_json, approval_status, approved_by, COALESCE(approved_at, '0001-01-01 00:00:00+00'::timestamptz), change_note, status, created_at, updated_at FROM notification_template_versions`
}

func scanTemplateVersions(rows pgx.Rows) ([]*biz.NotificationTemplateVersion, error) {
	var list []*biz.NotificationTemplateVersion
	for rows.Next() {
		item, err := scanTemplateVersion(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func scanTemplateVersion(row pgx.Row) (*biz.NotificationTemplateVersion, error) {
	item := &biz.NotificationTemplateVersion{}
	err := row.Scan(&item.ID, &item.TemplateID, &item.Version, &item.Channel, &item.Language, &item.Subject, &item.Body, &item.PayloadSchemaJSON, &item.ApprovalStatus, &item.ApprovedBy, &item.ApprovedAt, &item.ChangeNote, &item.Status, &item.CreatedAt, &item.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
}
