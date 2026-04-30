package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/arda-labs/arda/arda-be-go/services/notification-service/internal/biz"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var templatePlaceholderPattern = regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_.-]+)\s*\}\}`)

func (r *NotificationRepo) CreateNotificationRequest(ctx context.Context, item *biz.NotificationRequest) (*biz.NotificationRequest, error) {
	if existing, err := r.getNotificationRequestByIdempotencyKey(ctx, item.IdempotencyKey); err == nil {
		return existing, nil
	} else if !errors.Is(err, biz.ErrNotFound) {
		return nil, err
	}

	tx, err := r.data.db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `
		INSERT INTO notification_requests (
			source_service, event_type, correlation_id, idempotency_key, template_code,
			recipient_type, recipient_id, recipient_address, channels, language, payload_json, priority, status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id::text`,
		item.SourceService, item.EventType, item.CorrelationID, item.IdempotencyKey, item.TemplateCode,
		item.RecipientType, item.RecipientID, item.RecipientAddress, item.Channels, item.Language, item.PayloadJSON,
		item.Priority, item.Status).Scan(&item.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return r.getNotificationRequestByIdempotencyKey(ctx, item.IdempotencyKey)
		}
		return nil, err
	}

	payload, err := payloadMap(item.PayloadJSON)
	if err != nil {
		return nil, biz.ErrInvalidArgument
	}
	for _, channel := range item.Channels {
		version, err := getApprovedTemplateForDelivery(ctx, tx, item.TemplateCode, channel, item.Language)
		if err != nil {
			return nil, err
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO notification_deliveries (
				request_id, template_version_id, channel, recipient_type, recipient_id, recipient_address,
				subject, body, status, max_attempts, priority
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'QUEUED', 5, $9)`,
			item.ID, version.ID, channel, item.RecipientType, item.RecipientID, item.RecipientAddress,
			renderTemplate(version.Subject, payload), renderTemplate(version.Body, payload), item.Priority); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.GetNotificationRequest(ctx, item.ID)
}

func (r *NotificationRepo) GetNotificationRequest(ctx context.Context, id string) (*biz.NotificationRequest, error) {
	request, err := r.getNotificationRequest(ctx, id)
	if err != nil {
		return nil, err
	}
	deliveries, err := r.listDeliveriesByRequestID(ctx, id)
	if err != nil {
		return nil, err
	}
	request.Deliveries = deliveries
	return request, nil
}

func (r *NotificationRepo) getNotificationRequestByIdempotencyKey(ctx context.Context, key string) (*biz.NotificationRequest, error) {
	row := r.data.db.Pool.QueryRow(ctx, requestSelect()+` WHERE idempotency_key=$1`, key)
	request, err := scanNotificationRequest(row)
	if err != nil {
		return nil, err
	}
	deliveries, err := r.listDeliveriesByRequestID(ctx, request.ID)
	if err != nil {
		return nil, err
	}
	request.Deliveries = deliveries
	return request, nil
}

func (r *NotificationRepo) getNotificationRequest(ctx context.Context, id string) (*biz.NotificationRequest, error) {
	row := r.data.db.Pool.QueryRow(ctx, requestSelect()+` WHERE id=$1`, id)
	return scanNotificationRequest(row)
}

func getApprovedTemplateForDelivery(ctx context.Context, q pgx.Tx, templateCode, channel, language string) (*biz.NotificationTemplateVersion, error) {
	row := q.QueryRow(ctx, templateVersionSelect()+`
		WHERE template_id = (
			SELECT id FROM notification_templates WHERE code=$1 AND status='ACTIVE' AND deleted_at IS NULL
		)
		  AND channel=$2
		  AND language=$3
		  AND approval_status='APPROVED'
		  AND status='ACTIVE'
		  AND deleted_at IS NULL
		ORDER BY version DESC
		LIMIT 1`, templateCode, channel, language)
	return scanTemplateVersion(row)
}

func requestSelect() string {
	return `SELECT id::text, source_service, event_type, correlation_id, idempotency_key, template_code, recipient_type, recipient_id, recipient_address, channels, language, payload_json, priority, status, created_at, updated_at FROM notification_requests`
}

func scanNotificationRequest(row pgx.Row) (*biz.NotificationRequest, error) {
	item := &biz.NotificationRequest{}
	err := row.Scan(&item.ID, &item.SourceService, &item.EventType, &item.CorrelationID, &item.IdempotencyKey, &item.TemplateCode, &item.RecipientType, &item.RecipientID, &item.RecipientAddress, &item.Channels, &item.Language, &item.PayloadJSON, &item.Priority, &item.Status, &item.CreatedAt, &item.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
}

func payloadMap(raw string) (map[string]any, error) {
	out := map[string]any{}
	if strings.TrimSpace(raw) == "" {
		return out, nil
	}
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil, err
	}
	return out, nil
}

func renderTemplate(template string, payload map[string]any) string {
	return templatePlaceholderPattern.ReplaceAllStringFunc(template, func(match string) string {
		parts := templatePlaceholderPattern.FindStringSubmatch(match)
		if len(parts) != 2 {
			return match
		}
		if value, ok := payload[parts[1]]; ok {
			return fmt.Sprint(value)
		}
		return match
	})
}
