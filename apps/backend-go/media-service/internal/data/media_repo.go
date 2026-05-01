package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/arda-labs/arda/apps/backend-go/media-service/internal/biz"
)

type mediaRepo struct {
	data *Data
}

func NewMediaRepo(data *Data) biz.MediaRepo {
	return &mediaRepo{data: data}
}

func (r *mediaRepo) Create(ctx context.Context, media *biz.MediaMetadata) (*biz.MediaMetadata, error) {
	err := r.data.db.QueryRowContext(ctx, `
		INSERT INTO media_metadata (id, filename, content_type, size_bytes, bucket, object_key, owner_id, module, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at, updated_at`,
		media.ID, media.Filename, media.ContentType, media.SizeBytes, media.Bucket, media.ObjectKey, media.OwnerID, media.Module, media.Status,
	).Scan(&media.CreatedAt, &media.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("creating media metadata: %w", err)
	}
	return media, nil
}

func (r *mediaRepo) GetByID(ctx context.Context, id string) (*biz.MediaMetadata, error) {
	media := &biz.MediaMetadata{}
	err := r.data.db.QueryRowContext(ctx, `
		SELECT id, filename, content_type, size_bytes, bucket, object_key, owner_id, module, status, created_at, updated_at
		FROM media_metadata
		WHERE id = $1 AND status <> $2`, id, biz.StatusDeleted,
	).Scan(&media.ID, &media.Filename, &media.ContentType, &media.SizeBytes, &media.Bucket, &media.ObjectKey, &media.OwnerID, &media.Module, &media.Status, &media.CreatedAt, &media.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getting media metadata: %w", err)
	}
	return media, nil
}

func (r *mediaRepo) MarkReady(ctx context.Context, id string) (*biz.MediaMetadata, error) {
	media := &biz.MediaMetadata{}
	err := r.data.db.QueryRowContext(ctx, `
		UPDATE media_metadata
		SET status = $2, updated_at = now()
		WHERE id = $1 AND status = $3
		RETURNING id, filename, content_type, size_bytes, bucket, object_key, owner_id, module, status, created_at, updated_at`,
		id, biz.StatusReady, biz.StatusPending,
	).Scan(&media.ID, &media.Filename, &media.ContentType, &media.SizeBytes, &media.Bucket, &media.ObjectKey, &media.OwnerID, &media.Module, &media.Status, &media.CreatedAt, &media.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("marking media ready: %w", err)
	}
	return media, nil
}

func (r *mediaRepo) MarkDeleted(ctx context.Context, id string) error {
	_, err := r.data.db.ExecContext(ctx, `
		UPDATE media_metadata
		SET status = $2, updated_at = now()
		WHERE id = $1`, id, biz.StatusDeleted)
	if err != nil {
		return fmt.Errorf("marking media deleted: %w", err)
	}
	return nil
}

func (r *mediaRepo) ListPendingBefore(ctx context.Context, before time.Time, limit int) ([]*biz.MediaMetadata, error) {
	rows, err := r.data.db.QueryContext(ctx, `
		SELECT id, filename, content_type, size_bytes, bucket, object_key, owner_id, module, status, created_at, updated_at
		FROM media_metadata
		WHERE status = $1 AND created_at < $2
		ORDER BY created_at ASC
		LIMIT $3`, biz.StatusPending, before, limit)
	if err != nil {
		return nil, fmt.Errorf("listing pending media: %w", err)
	}
	defer rows.Close()

	items := make([]*biz.MediaMetadata, 0)
	for rows.Next() {
		media := &biz.MediaMetadata{}
		if err := rows.Scan(&media.ID, &media.Filename, &media.ContentType, &media.SizeBytes, &media.Bucket, &media.ObjectKey, &media.OwnerID, &media.Module, &media.Status, &media.CreatedAt, &media.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, media)
	}
	return items, rows.Err()
}
