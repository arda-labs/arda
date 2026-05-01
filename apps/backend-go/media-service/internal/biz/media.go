package biz

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	StatusPending = "PENDING"
	StatusReady   = "READY"
	StatusDeleted = "DELETED"
)

type MediaMetadata struct {
	ID          string    `json:"id"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"content_type"`
	SizeBytes   int64     `json:"size_bytes"`
	Bucket      string    `json:"bucket"`
	ObjectKey   string    `json:"object_key"`
	OwnerID     string    `json:"owner_id"`
	Module      string    `json:"module"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type InitUploadInput struct {
	Filename    string
	ContentType string
	SizeBytes   int64
	OwnerID     string
	Module      string
}

type InitUploadResult struct {
	Media     *MediaMetadata
	UploadURL string
	ExpiresAt time.Time
}

type DownloadURLResult struct {
	Media       *MediaMetadata
	DownloadURL string
	ExpiresAt   time.Time
}

type MediaRepo interface {
	Create(ctx context.Context, media *MediaMetadata) (*MediaMetadata, error)
	GetByID(ctx context.Context, id string) (*MediaMetadata, error)
	MarkReady(ctx context.Context, id string) (*MediaMetadata, error)
	MarkDeleted(ctx context.Context, id string) error
	ListPendingBefore(ctx context.Context, before time.Time, limit int) ([]*MediaMetadata, error)
}

type StorageRepo interface {
	PresignPutObject(ctx context.Context, bucket, objectKey, contentType string, ttl time.Duration) (string, time.Time, error)
	PresignGetObject(ctx context.Context, bucket, objectKey string, ttl time.Duration) (string, time.Time, error)
	HeadObject(ctx context.Context, bucket, objectKey string) (int64, error)
	DeleteObject(ctx context.Context, bucket, objectKey string) error
}

type MediaUsecase struct {
	repo           MediaRepo
	storage        StorageRepo
	bucket         string
	uploadURLTTL   time.Duration
	downloadURLTTL time.Duration
}

func NewMediaUsecase(repo MediaRepo, storage StorageRepo, bucket string, uploadURLTTL, downloadURLTTL time.Duration) *MediaUsecase {
	return &MediaUsecase{repo: repo, storage: storage, bucket: bucket, uploadURLTTL: uploadURLTTL, downloadURLTTL: downloadURLTTL}
}

func (uc *MediaUsecase) InitUpload(ctx context.Context, in InitUploadInput) (*InitUploadResult, error) {
	if strings.TrimSpace(in.Filename) == "" {
		return nil, fmt.Errorf("filename is required")
	}
	if strings.TrimSpace(in.ContentType) == "" {
		return nil, fmt.Errorf("content_type is required")
	}
	if in.SizeBytes <= 0 {
		return nil, fmt.Errorf("size_bytes must be greater than zero")
	}
	if strings.TrimSpace(in.OwnerID) == "" {
		return nil, fmt.Errorf("owner_id is required")
	}

	id := uuid.NewString()
	module := normalizeModule(in.Module)
	now := time.Now().UTC()
	objectKey := buildObjectKey(module, now, id, in.Filename)
	media := &MediaMetadata{
		ID:          id,
		Filename:    in.Filename,
		ContentType: in.ContentType,
		SizeBytes:   in.SizeBytes,
		Bucket:      uc.bucket,
		ObjectKey:   objectKey,
		OwnerID:     in.OwnerID,
		Module:      module,
		Status:      StatusPending,
	}

	created, err := uc.repo.Create(ctx, media)
	if err != nil {
		return nil, err
	}
	url, expiresAt, err := uc.storage.PresignPutObject(ctx, created.Bucket, created.ObjectKey, created.ContentType, uc.uploadURLTTL)
	if err != nil {
		return nil, err
	}
	return &InitUploadResult{Media: created, UploadURL: url, ExpiresAt: expiresAt}, nil
}

func (uc *MediaUsecase) ConfirmUpload(ctx context.Context, id string) (*MediaMetadata, error) {
	media, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if media == nil {
		return nil, fmt.Errorf("media not found")
	}
	if media.Status != StatusPending {
		return nil, fmt.Errorf("media is not pending")
	}
	size, err := uc.storage.HeadObject(ctx, media.Bucket, media.ObjectKey)
	if err != nil {
		return nil, err
	}
	if size != media.SizeBytes {
		return nil, fmt.Errorf("uploaded size mismatch: expected %d, got %d", media.SizeBytes, size)
	}
	return uc.repo.MarkReady(ctx, id)
}

func (uc *MediaUsecase) GetDownloadURL(ctx context.Context, id, ownerID string) (*DownloadURLResult, error) {
	media, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if media == nil {
		return nil, fmt.Errorf("media not found")
	}
	if media.Status != StatusReady {
		return nil, fmt.Errorf("media is not ready")
	}
	if ownerID != "" && ownerID != media.OwnerID {
		return nil, fmt.Errorf("media access denied")
	}
	url, expiresAt, err := uc.storage.PresignGetObject(ctx, media.Bucket, media.ObjectKey, uc.downloadURLTTL)
	if err != nil {
		return nil, err
	}
	return &DownloadURLResult{Media: media, DownloadURL: url, ExpiresAt: expiresAt}, nil
}

func (uc *MediaUsecase) CleanupPending(ctx context.Context, before time.Time, limit int) (int, error) {
	items, err := uc.repo.ListPendingBefore(ctx, before, limit)
	if err != nil {
		return 0, err
	}
	cleaned := 0
	for _, item := range items {
		if err := uc.storage.DeleteObject(ctx, item.Bucket, item.ObjectKey); err != nil {
			return cleaned, err
		}
		if err := uc.repo.MarkDeleted(ctx, item.ID); err != nil {
			return cleaned, err
		}
		cleaned++
	}
	return cleaned, nil
}

var modulePattern = regexp.MustCompile(`[^a-z0-9_-]+`)

func normalizeModule(module string) string {
	module = strings.ToLower(strings.TrimSpace(module))
	module = modulePattern.ReplaceAllString(module, "-")
	module = strings.Trim(module, "-")
	if module == "" {
		return "general"
	}
	return module
}

func buildObjectKey(module string, now time.Time, id, filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if len(ext) > 16 {
		ext = ""
	}
	return fmt.Sprintf("%s/%04d/%02d/%s%s", module, now.Year(), now.Month(), id, ext)
}
