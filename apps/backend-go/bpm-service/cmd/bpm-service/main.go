package main

import (
	"context"
	"log"
	"os"
	"github.com/arda-labs/arda/apps/backend-go/bpm-service/internal/data/events"
	"github.com/arda-labs/arda/apps/backend-go/bpm-service/internal/worker"
	"github.com/go-kratos/kratos/v2"
)

func main() {
	log.Println("Starting BPM Service (Go) with Kafka Event Bus and Zeebe Integration...")

	// Default to thinkcenter if running on host, but keep localhost for local docker setup
	// In production, these will be injected via K8s Environment Variables
	brokers := []string{"thinkcenter:9092"}
	zeebeAddr := "thinkcenter:26500"

	// Check if running on windows/localhost
	if os.Getenv("ENV") == "LOCAL" {
		brokers = []string{"localhost:9092"}
		zeebeAddr = "localhost:26500"
	}

	// 1. Initialize Publisher
	publisher := events.NewKafkaPublisher(brokers)

	// 2. Initialize Zeebe Client
	zeebe, err := worker.NewZeebeClient(zeebeAddr)
	if err != nil {
		log.Fatalf("Failed to initialize Zeebe: %v", err)
	}

	// 3. Initialize Generic Worker (with Publisher)
	_ = worker.NewGenericWorker(publisher)

	// 4. Initialize & Start CRM Event Consumer (with Zeebe)
	crmConsumer := worker.NewCRMEventConsumer(brokers, zeebe)
	go crmConsumer.Start(context.Background())

	app := kratos.New(
		kratos.Name("bpm-service"),
		kratos.Version("2.1.0"),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
