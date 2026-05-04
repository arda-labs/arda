package worker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/arda-labs/arda/arda-be-go/services/bpm-service/internal/biz"
	"github.com/arda-labs/arda/arda-be-go/services/bpm-service/internal/data/events"
)

type CRMEventConsumer struct {
	subscriber *events.KafkaSubscriber
	zeebe      *ZeebeClient
	instUC     *biz.InstanceUseCase
	defUC      *biz.DefinitionUseCase
	eventUC    *biz.EventUseCase
}

func NewCRMEventConsumer(brokers []string, zeebe *ZeebeClient, instUC *biz.InstanceUseCase, defUC *biz.DefinitionUseCase, eventUC *biz.EventUseCase) *CRMEventConsumer {
	return &CRMEventConsumer{
		subscriber: events.NewKafkaSubscriber(brokers, "crm-events", "bpm-group"),
		zeebe:      zeebe,
		instUC:     instUC,
		defUC:      defUC,
		eventUC:    eventUC,
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
			zeebeKey, err := c.zeebe.StartInstance(ctx, "process_customer_registration", event.Data)
			if err != nil {
				log.Printf("BPM: Failed to start Zeebe process: %v\n", err)
				return err
			}

			// Look up the process definition to get its ID
			defs, _, err := c.defUC.List(ctx, biz.DefinitionFilter{Keyword: "process_customer_registration", PageSize: 1})
			if err != nil || len(defs) == 0 {
				log.Printf("BPM: Definition not found for process_customer_registration, skipping DB persist")
				return nil
			}

			// Persist instance to PostgreSQL
			varsJSON, _ := json.Marshal(event.Data)
			inst := &biz.ProcessInstance{
				ZeebeInstanceKey:   zeebeKey,
				ProcessDefinitionID: defs[0].ID,
				Status:             "ACTIVE",
				CurrentStep:        "start",
				Variables:          string(varsJSON),
				SLAStatus:          "ON_TRACK",
			}
			created, err := c.instUC.Create(ctx, inst)
			if err != nil {
				log.Printf("BPM: Failed to persist instance: %v\n", err)
				return err
			}

			// Record the start event
			_, _ = c.eventUC.Create(ctx, &biz.ProcessEvent{
				ProcessInstanceID: created.ID,
				EventType:         "INSTANCE_CREATED",
				Source:            "crm-consumer",
				Data:              string(varsJSON),
			})

			log.Printf("BPM: Persisted instance %s (Zeebe key: %d)\n", created.ID, zeebeKey)
		}

		return nil
	})
}
