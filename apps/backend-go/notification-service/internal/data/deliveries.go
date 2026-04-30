package data

import (
	"context"
	"errors"

	"github.com/arda-labs/arda/arda-be-go/pkg/pagination"
	"github.com/arda-labs/arda/arda-be-go/services/notification-service/internal/biz"
	"github.com/jackc/pgx/v5"
)

func (r *NotificationRepo) ListDeliveries(ctx context.Context, filter biz.DeliveryFilter) ([]*biz.NotificationDelivery, string, error) {
	page := pagination.Normalize(filter.PageSize, filter.PageToken)
	rows, err := r.data.db.Pool.Query(ctx, deliverySelect()+`
		WHERE ($1 = '' OR status = $1)
		  AND ($2 = '' OR channel = $2)
		  AND ($3 = '' OR recipient_id = $3)
		ORDER BY created_at DESC
		LIMIT $4 OFFSET $5`, filter.Status, filter.Channel, filter.RecipientID, page.Limit+1, page.Offset)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()
	list, err := scanDeliveries(rows)
	if err != nil {
		return nil, "", err
	}
	next := pagination.NextOffsetToken(len(list), page.Limit, page.Offset)
	if len(list) > page.Limit {
		list = list[:page.Limit]
	}
	return list, next, nil
}

func (r *NotificationRepo) RetryDelivery(ctx context.Context, id, actor string) (*biz.NotificationDelivery, error) {
	_ = actor
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE notification_deliveries
		SET status='QUEUED', next_attempt_at=now(), locked_by='', locked_at=NULL, error_message='', updated_at=now()
		WHERE id=$1 AND status IN ('FAILED', 'RETRYING', 'DEAD_LETTER')`, id)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.getDelivery(ctx, id)
}

func (r *NotificationRepo) ClaimDueDeliveries(ctx context.Context, workerID string, limit int) ([]*biz.NotificationDelivery, error) {
	rows, err := r.data.db.Pool.Query(ctx, `
		WITH due AS (
			SELECT id
			FROM notification_deliveries
			WHERE status IN ('QUEUED', 'RETRYING')
			  AND next_attempt_at <= now()
			ORDER BY priority ASC, created_at ASC
			LIMIT $2
			FOR UPDATE SKIP LOCKED
		)
		UPDATE notification_deliveries d
		SET status='CLAIMED', locked_by=$1, locked_at=now(), attempt_count=d.attempt_count+1, updated_at=now()
		FROM due
		WHERE d.id = due.id
		RETURNING d.id::text, d.request_id::text, d.template_version_id::text, d.channel, d.recipient_type, d.recipient_id, d.recipient_address,
		          d.subject, d.body, d.status, d.attempt_count, d.max_attempts, d.next_attempt_at, d.locked_by,
		          COALESCE(d.locked_at, '0001-01-01 00:00:00+00'::timestamptz), d.provider_code, d.provider_message_id,
		          d.provider_response_json, d.error_message, d.priority, d.created_at, d.updated_at`, workerID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDeliveries(rows)
}

func (r *NotificationRepo) MarkDeliveryDelivered(ctx context.Context, id, providerCode, providerMessageID, providerResponseJSON string) (*biz.NotificationDelivery, error) {
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE notification_deliveries
		SET status='DELIVERED', provider_code=$2, provider_message_id=$3, provider_response_json=$4,
		    locked_by='', locked_at=NULL, error_message='', updated_at=now()
		WHERE id=$1`, id, providerCode, providerMessageID, providerResponseJSON)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.getDelivery(ctx, id)
}

func (r *NotificationRepo) MarkDeliveryFailed(ctx context.Context, id, message string, retryAfterSeconds int) (*biz.NotificationDelivery, error) {
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE notification_deliveries
		SET status = CASE WHEN attempt_count >= max_attempts THEN 'DEAD_LETTER' ELSE 'RETRYING' END,
		    next_attempt_at = CASE WHEN attempt_count >= max_attempts THEN next_attempt_at ELSE now() + ($2::text || ' seconds')::interval END,
		    locked_by='', locked_at=NULL, error_message=$3, updated_at=now()
		WHERE id=$1`, id, retryAfterSeconds, message)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.getDelivery(ctx, id)
}

func (r *NotificationRepo) listDeliveriesByRequestID(ctx context.Context, requestID string) ([]*biz.NotificationDelivery, error) {
	rows, err := r.data.db.Pool.Query(ctx, deliverySelect()+` WHERE request_id=$1 ORDER BY created_at ASC`, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDeliveries(rows)
}

func (r *NotificationRepo) getDelivery(ctx context.Context, id string) (*biz.NotificationDelivery, error) {
	row := r.data.db.Pool.QueryRow(ctx, deliverySelect()+` WHERE id=$1`, id)
	return scanDelivery(row)
}

func deliverySelect() string {
	return `SELECT id::text, request_id::text, template_version_id::text, channel, recipient_type, recipient_id, recipient_address, subject, body, status, attempt_count, max_attempts, next_attempt_at, locked_by, COALESCE(locked_at, '0001-01-01 00:00:00+00'::timestamptz), provider_code, provider_message_id, provider_response_json, error_message, priority, created_at, updated_at FROM notification_deliveries`
}

func scanDeliveries(rows pgx.Rows) ([]*biz.NotificationDelivery, error) {
	var list []*biz.NotificationDelivery
	for rows.Next() {
		item, err := scanDelivery(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func scanDelivery(row pgx.Row) (*biz.NotificationDelivery, error) {
	item := &biz.NotificationDelivery{}
	err := row.Scan(&item.ID, &item.RequestID, &item.TemplateVersionID, &item.Channel, &item.RecipientType, &item.RecipientID, &item.RecipientAddress, &item.Subject, &item.Body, &item.Status, &item.AttemptCount, &item.MaxAttempts, &item.NextAttemptAt, &item.LockedBy, &item.LockedAt, &item.ProviderCode, &item.ProviderMessageID, &item.ProviderResponseJSON, &item.ErrorMessage, &item.Priority, &item.CreatedAt, &item.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
}
