package main

import (
	"context"
	"flag"
	"os"

	"github.com/arda-labs/arda/arda-be-go/services/bpm-service/internal/conf"
	"github.com/arda-labs/arda/arda-be-go/services/bpm-service/internal/data/events"
	"github.com/arda-labs/arda/arda-be-go/services/bpm-service/internal/service"
	"github.com/arda-labs/arda/arda-be-go/services/bpm-service/internal/worker"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	Name    string
	Version string
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
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

	// Wire dependencies (DB, repos, use cases, servers)
	app, cleanup, bpmService, err := wireApp(bc.Server, bc.Data, bc.Jwt, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// Initialize background workers (Zeebe client, Kafka consumers)
	initWorkers(bc.Data, bpmService)

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func initWorkers(data *conf.Data, bpmService *service.BPMService) {
	if data == nil {
		return
	}

	brokers := data.Kafka.GetBrokers()
	if len(brokers) == 0 {
		brokers = []string{"thinkcenter:9092"}
	}

	zeebeAddr := data.Zeebe.GetAddr()
	if zeebeAddr == "" {
		zeebeAddr = "thinkcenter:26500"
	}

	// Override with LOCAL env for local development
	if os.Getenv("ENV") == "LOCAL" {
		brokers = []string{"localhost:9092"}
		zeebeAddr = "localhost:26500"
	}

	// Initialize Publisher
	publisher := events.NewKafkaPublisher(brokers)

	// Initialize Zeebe Client
	zeebe, err := worker.NewZeebeClient(zeebeAddr)
	if err != nil {
		log.NewHelper(log.DefaultLogger).Warnf("Failed to initialize Zeebe: %v", err)
		return
	}

	// Wire Zeebe deployer into BPMService
	bpmService.SetDeployer(zeebe.DeployBPMN)

	// Initialize Generic Worker (with Publisher)
	_ = worker.NewGenericWorker(publisher)

	// Initialize & Start CRM Event Consumer (with Zeebe, instUC, defUC, eventUC)
	defUC, instUC, eventUC, _ := bpmService.UseCases()
	crmConsumer := worker.NewCRMEventConsumer(brokers, zeebe, instUC, defUC, eventUC)
	go crmConsumer.Start(context.Background())

}
