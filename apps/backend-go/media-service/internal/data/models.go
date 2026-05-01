package data

import "time"

type MediaMetadata struct {
	ID          string
	Filename    string
	ContentType string
	SizeBytes   int64
	S3Key       string
	OwnerID     *string
	CreatedAt   time.Time
}
