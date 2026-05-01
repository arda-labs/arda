package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arda-labs/arda/apps/backend-go/media-service/internal/biz"
	"github.com/arda-labs/arda/apps/backend-go/media-service/internal/conf"
	"github.com/arda-labs/arda/apps/backend-go/media-service/internal/data"
	"github.com/arda-labs/arda/apps/backend-go/media-service/internal/server"
	"github.com/arda-labs/arda/apps/backend-go/media-service/internal/service"
	"github.com/arda-labs/arda/apps/backend-go/media-service/internal/storage"
	"github.com/arda-labs/arda/apps/backend-go/media-service/internal/worker"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
)

var (
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dataStore, cleanup, err := data.NewData(bc.Data.Database.Source)
	if err != nil {
		log.Fatalf("init data: %v", err)
	}
	defer cleanup()

	storageRepo, err := storage.NewS3Storage(ctx, storage.S3Config{
		Endpoint:        bc.Storage.S3.Endpoint,
		Region:          bc.Storage.S3.Region,
		AccessKeyID:     bc.Storage.S3.AccessKey,
		SecretAccessKey: bc.Storage.S3.SecretKey,
	})
	if err != nil {
		log.Fatalf("init s3 storage: %v", err)
	}

	uploadTTL, _ := time.ParseDuration(bc.Storage.S3.UploadUrlTtl)
	if uploadTTL == 0 {
		uploadTTL = 15 * time.Minute
	}
	downloadTTL, _ := time.ParseDuration(bc.Storage.S3.DownloadUrlTtl)
	if downloadTTL == 0 {
		downloadTTL = 15 * time.Minute
	}

	mediaRepo := data.NewMediaRepo(dataStore)
	mediaUsecase := biz.NewMediaUsecase(mediaRepo, storageRepo, bc.Storage.S3.Bucket, uploadTTL, downloadTTL)
	mediaService := service.NewMediaService(mediaUsecase)

	cleanupWorker := worker.NewCleanupWorker(mediaUsecase, time.Hour, 24*time.Hour, 100)
	go cleanupWorker.Start(ctx)

	srv := &http.Server{Addr: bc.Server.Http.Addr, Handler: server.NewHTTPServer(mediaService)}
	go func() {
		log.Printf("Media Service starting on %s", bc.Server.Http.Addr)
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
