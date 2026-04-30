package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/arda-labs/arda/arda-be-go/services/notification-service/internal/biz"
	"github.com/arda-labs/arda/arda-be-go/services/notification-service/internal/conf"
	"github.com/arda-labs/arda/arda-be-go/services/notification-service/internal/data"
	"github.com/arda-labs/arda/arda-be-go/services/notification-service/internal/server"
	"github.com/arda-labs/arda/arda-be-go/services/notification-service/internal/service"
	"github.com/arda-labs/arda/arda-be-go/services/notification-service/internal/worker"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	_ "go.uber.org/automaxprocs"
)

var (
	Name     = "notification-service"
	Version  string
	flagconf string
	id, _    = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "configs/config.yaml", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()
	if !filepath.IsAbs(flagconf) {
		if cwd, err := os.Getwd(); err == nil {
			path := filepath.Join(cwd, flagconf)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				path = filepath.Join(filepath.Dir(cwd), flagconf)
				if _, err := os.Stat(path); os.IsNotExist(err) {
					path = filepath.Join(filepath.Dir(filepath.Dir(cwd)), flagconf)
				}
			}
			flagconf = path
		}
	}

	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	c := config.New(config.WithSource(file.NewSource(flagconf), env.NewSource("")))
	defer c.Close()
	if err := c.Load(); err != nil {
		panic(err)
	}
	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	if envDB := os.Getenv("DATABASE_URL"); envDB != "" {
		bc.Data.Database.Source = envDB
	}

	d, cleanup, err := data.NewData(bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()
	repo := data.NewNotificationRepo(d)
	uc := biz.NewNotificationUsecase(repo)
	svc := service.NewNotificationService(uc)
	hs := server.NewHTTPServer(bc.Server, bc.Jwt, svc, logger)
	gs := server.NewGRPCServer(bc.Server, svc, logger)
	workerCtx, stopWorker := context.WithCancel(context.Background())
	defer stopWorker()
	worker.NewDeliveryWorker(uc, logger, id, 5*time.Second, 20).Start(workerCtx)

	app := kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Logger(logger),
		kratos.Server(gs, hs),
	)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
