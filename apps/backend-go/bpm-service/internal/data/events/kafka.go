package events

import (
	"context"
	"encoding/json"
	"log"
	"github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(brokers []string) *KafkaPublisher {
	return &KafkaPublisher{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *KafkaPublisher) Publish(ctx context.Context, topic string, event interface{}) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Value: payload,
	})

	if err != nil {
		log.Printf("Kafka: Failed to publish message to %s: %v\n", topic, err)
		return err
	}

	log.Printf("Kafka: Published message to %s\n", topic)
	return nil
}

type KafkaSubscriber struct {
	reader *kafka.Reader
}

func NewKafkaSubscriber(brokers []string, topic, groupID string) *KafkaSubscriber {
	return &KafkaSubscriber{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  groupID,
			Topic:    topic,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
	}
}

func (s *KafkaSubscriber) Consume(ctx context.Context, handler func([]byte) error) {
	for {
		m, err := s.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Kafka: Error reading message: %v\n", err)
			break
		}
		
		log.Printf("Kafka: Received message from %s\n", m.Topic)
		if err := handler(m.Value); err != nil {
			log.Printf("Kafka: Error handling message: %v\n", err)
		}
	}
}
