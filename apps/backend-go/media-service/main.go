package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arda-labs/arda/apps/backend-go/media-service/internal/biz"
	"github.com/arda-labs/arda/apps/backend-go/media-service/internal/data"
	"github.com/arda-labs/arda/apps/backend-go/media-service/internal/server"
	"github.com/arda-labs/arda/apps/backend-go/media-service/internal/service"
	"github.com/arda-labs/arda/apps/backend-go/media-service/internal/storage"
	"github.com/arda-labs/arda/apps/backend-go/media-service/internal/worker"
)

type config struct {
	Port              string
	DatabaseURL       string
	S3Endpoint        string
	S3Region          string
	S3Bucket          string
	S3AccessKeyID     string
	S3SecretAccessKey string
	UploadURLTTL      time.Duration
	DownloadURLTTL    time.Duration
}

func main() {
	cfg := loadConfig()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dataStore, cleanup, err := data.NewData(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("init data: %v", err)
	}
	defer cleanup()

	storageRepo, err := storage.NewS3Storage(ctx, storage.S3Config{
		Endpoint:        cfg.S3Endpoint,
		Region:          cfg.S3Region,
		AccessKeyID:     cfg.S3AccessKeyID,
		SecretAccessKey: cfg.S3SecretAccessKey,
	})
	if err != nil {
		log.Fatalf("init s3 storage: %v", err)
	}

	mediaRepo := data.NewMediaRepo(dataStore)
	mediaUsecase := biz.NewMediaUsecase(mediaRepo, storageRepo, cfg.S3Bucket, cfg.UploadURLTTL, cfg.DownloadURLTTL)
	mediaService := service.NewMediaService(mediaUsecase)

	cleanupWorker := worker.NewCleanupWorker(mediaUsecase, time.Hour, 24*time.Hour, 100)
	go cleanupWorker.Start(ctx)

	addr := ":" + cfg.Port
	srv := &http.Server{Addr: addr, Handler: server.NewHTTPServer(mediaService)}
	go func() {
		log.Printf("Media Service starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("http shutdown failed: %v", err)
	}
}

func loadConfig() config {
	return config{
		Port:              env("PORT", "8080"),
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		S3Endpoint:        requiredEnv("STORAGE_S3_ENDPOINT"),
		S3Region:          env("STORAGE_S3_REGION", "us-east-1"),
		S3Bucket:          env("STORAGE_S3_BUCKET", "arda-media"),
		S3AccessKeyID:     env("STORAGE_S3_ACCESS_KEY_ID", "admin"),
		S3SecretAccessKey: env("STORAGE_S3_SECRET_ACCESS_KEY", "admin"),
		UploadURLTTL:      durationEnv("UPLOAD_URL_TTL", 15*time.Minute),
		DownloadURLTTL:    durationEnv("DOWNLOAD_URL_TTL", 15*time.Minute),
	}
}

func env(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func requiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s is required", key)
	}
	return value
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("invalid %s: %v", key, err)
	}
	return duration
}

func (c config) String() string {
	return fmt.Sprintf("port=%s endpoint=%s bucket=%s", c.Port, c.S3Endpoint, c.S3Bucket)
}
