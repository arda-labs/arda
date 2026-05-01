package data

import "time"

type MediaMetadata struct {
	ID          string
	Filename    string
	ContentType string
	SizeBytes   int64
	Bucket      string
	ObjectKey   string
	OwnerID     string
	Module      string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
