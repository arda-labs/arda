package storage

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Config struct {
	Endpoint        string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
}

type S3Storage struct {
	client    *s3.Client
	presigner *s3.PresignClient
}

func NewS3Client(ctx context.Context, cfg S3Config) (*s3.Client, error) {
	awsCfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, "")),
	)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(awsCfg, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(cfg.Endpoint)
		options.UsePathStyle = true
	}), nil
}

func NewS3Storage(ctx context.Context, cfg S3Config) (*S3Storage, error) {
	client, err := NewS3Client(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &S3Storage{client: client, presigner: s3.NewPresignClient(client)}, nil
}

func (s *S3Storage) PresignPutObject(ctx context.Context, bucket, objectKey, contentType string, ttl time.Duration) (string, time.Time, error) {
	expiresAt := time.Now().UTC().Add(ttl)
	out, err := s.presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(objectKey),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", time.Time{}, err
	}
	return out.URL, expiresAt, nil
}

func (s *S3Storage) PresignGetObject(ctx context.Context, bucket, objectKey string, ttl time.Duration) (string, time.Time, error) {
	expiresAt := time.Now().UTC().Add(ttl)
	out, err := s.presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", time.Time{}, err
	}
	return out.URL, expiresAt, nil
}

func (s *S3Storage) HeadObject(ctx context.Context, bucket, objectKey string) (int64, error) {
	out, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return 0, err
	}
	if out.ContentLength == nil {
		return 0, nil
	}
	return *out.ContentLength, nil
}

func (s *S3Storage) DeleteObject(ctx context.Context, bucket, objectKey string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	})
	return err
}
