package arda.messaging.producer;

import arda.common.context.ArdaContext;
import arda.messaging.event.ArdaEvent;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Component;

import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

@Component
public class ArdaEventProducer {
    private final KafkaTemplate<String, String> kafkaTemplate;
    private final ObjectMapper objectMapper;

    public ArdaEventProducer(KafkaTemplate<String, String> kafkaTemplate, ObjectMapper objectMapper) {
        this.kafkaTemplate = kafkaTemplate;
        this.objectMapper = objectMapper;
    }

    /**
     * Send an event synchronously.
     * Works perfectly with Virtual Threads.
     */
    public <T> void send(String topic, String eventType, T data, String serviceName) {
        ArdaContext ctx = ArdaContext.current();

        ArdaEvent<T> event = new ArdaEvent<>(
            serviceName,
            eventType,
            data,
            ctx.traceId(),
            ctx.userId()
        );

        try {
            String payload = objectMapper.writeValueAsString(event);
            // Block and wait for Kafka ACK.
            // In Virtual Threads, this won't block the OS thread.
            kafkaTemplate.send(topic, event.id(), payload)
                .get(5, TimeUnit.SECONDS);
        } catch (JsonProcessingException | InterruptedException | ExecutionException | TimeoutException e) {
            throw new RuntimeException("Failed to send Kafka event: " + eventType, e);
        }
    }
}
