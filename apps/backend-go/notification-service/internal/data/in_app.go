package data

import (
	"context"
	"errors"

	"github.com/arda-labs/arda/arda-be-go/pkg/pagination"
	"github.com/arda-labs/arda/arda-be-go/services/notification-service/internal/biz"
	"github.com/jackc/pgx/v5"
)

func (r *NotificationRepo) CreateInAppNotification(ctx context.Context, delivery *biz.NotificationDelivery) (*biz.InAppNotification, error) {
	var id string
	err := r.data.db.Pool.QueryRow(ctx, `
		INSERT INTO in_app_notifications (delivery_id, recipient_type, recipient_id, title, body, data_json, status)
		VALUES ($1, $2, $3, $4, $5, $6, 'UNREAD')
		ON CONFLICT (delivery_id) DO UPDATE
		SET title=EXCLUDED.title
		RETURNING id::text`,
		delivery.ID, delivery.RecipientType, delivery.RecipientID, delivery.Subject, delivery.Body, delivery.ProviderResponseJSON).Scan(&id)
	if err != nil {
		return nil, err
	}
	return r.getInAppNotification(ctx, id)
}

func (r *NotificationRepo) ListInAppNotifications(ctx context.Context, filter biz.InAppFilter) ([]*biz.InAppNotification, string, error) {
	page := pagination.Normalize(filter.PageSize, filter.PageToken)
	rows, err := r.data.db.Pool.Query(ctx, inAppSelect()+`
		WHERE recipient_type=$1
		  AND recipient_id=$2
		  AND ($3 = '' OR status = $3)
		ORDER BY created_at DESC
		LIMIT $4 OFFSET $5`, filter.RecipientType, filter.RecipientID, filter.Status, page.Limit+1, page.Offset)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()
	list, err := scanInAppNotifications(rows)
	if err != nil {
		return nil, "", err
	}
	next := pagination.NextOffsetToken(len(list), page.Limit, page.Offset)
	if len(list) > page.Limit {
		list = list[:page.Limit]
	}
	return list, next, nil
}

func (r *NotificationRepo) MarkInAppNotificationRead(ctx context.Context, id, actor string) (*biz.InAppNotification, error) {
	_ = actor
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE in_app_notifications
		SET status='READ', read_at=COALESCE(read_at, now())
		WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.getInAppNotification(ctx, id)
}

func (r *NotificationRepo) CountUnreadInAppNotifications(ctx context.Context, recipientType, recipientID string) (int, error) {
	var count int
	err := r.data.db.Pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM in_app_notifications
		WHERE recipient_type=$1
		  AND recipient_id=$2
		  AND status='UNREAD'`, recipientType, recipientID).Scan(&count)
	return count, err
}

func (r *NotificationRepo) MarkAllInAppNotificationsRead(ctx context.Context, recipientType, recipientID, actor string) (int, error) {
	_ = actor
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE in_app_notifications
		SET status='READ', read_at=COALESCE(read_at, now())
		WHERE recipient_type=$1
		  AND recipient_id=$2
		  AND status='UNREAD'`, recipientType, recipientID)
	if err != nil {
		return 0, err
	}
	return int(tag.RowsAffected()), nil
}

func (r *NotificationRepo) getInAppNotification(ctx context.Context, id string) (*biz.InAppNotification, error) {
	row := r.data.db.Pool.QueryRow(ctx, inAppSelect()+` WHERE id=$1`, id)
	return scanInAppNotification(row)
}

func inAppSelect() string {
	return `SELECT id::text, delivery_id::text, recipient_type, recipient_id, title, body, data_json, status, COALESCE(read_at, '0001-01-01 00:00:00+00'::timestamptz), created_at FROM in_app_notifications`
}

func scanInAppNotifications(rows pgx.Rows) ([]*biz.InAppNotification, error) {
	var list []*biz.InAppNotification
	for rows.Next() {
		item, err := scanInAppNotification(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func scanInAppNotification(row pgx.Row) (*biz.InAppNotification, error) {
	item := &biz.InAppNotification{}
	err := row.Scan(&item.ID, &item.DeliveryID, &item.RecipientType, &item.RecipientID, &item.Title, &item.Body, &item.DataJSON, &item.Status, &item.ReadAt, &item.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
}
