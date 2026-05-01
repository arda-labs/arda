package worker

import (
	"context"
	"encoding/json"
	"log"
	"github.com/arda-labs/arda/apps/backend-go/bpm-service/internal/data/events"
)

type CRMEventConsumer struct {
	subscriber *events.KafkaSubscriber
	zeebe      *ZeebeClient
}

func NewCRMEventConsumer(brokers []string, zeebe *ZeebeClient) *CRMEventConsumer {
	return &CRMEventConsumer{
		subscriber: events.NewKafkaSubscriber(brokers, "crm-events", "bpm-group"),
		zeebe:      zeebe,
	}
}

func (c *CRMEventConsumer) Start(ctx context.Context) {
	log.Println("BPM: Starting to consume CRM events from Kafka...")
	c.subscriber.Consume(ctx, func(data []byte) error {
		var event struct {
			Type string                 `json:"type"`
			Data map[string]interface{} `json:"data"`
		}
		if err := json.Unmarshal(data, &event); err != nil {
			return err
		}

		log.Printf("BPM: Received event %v from CRM\n", event.Type)

		if event.Type == "CUSTOMER_REGISTRATION_CREATED" {
			log.Println("BPM: Auto-starting Workflow for New Customer Registration...")
			_, err := c.zeebe.StartInstance(ctx, "process_customer_registration", event.Data)
			if err != nil {
				log.Printf("BPM: Failed to start Zeebe process: %v\n", err)
				return err
			}
		}

		return nil
	})
}
