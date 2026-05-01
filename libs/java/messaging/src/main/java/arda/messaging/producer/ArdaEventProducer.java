package arda.messaging.producer;

import arda.common.context.ArdaContext;
import arda.messaging.event.ArdaEvent;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Component;
import reactor.core.publisher.Mono;

@Component
public class ArdaEventProducer {
    private final KafkaTemplate<String, String> kafkaTemplate;
    private final ObjectMapper objectMapper;

    public ArdaEventProducer(KafkaTemplate<String, String> kafkaTemplate, ObjectMapper objectMapper) {
        this.kafkaTemplate = kafkaTemplate;
        this.objectMapper = objectMapper;
    }

    /**
     * Send an event with automatic Context propagation (traceId, userId).
     */
    public <T> Mono<Void> send(String topic, String eventType, T data, String serviceName) {
        return ArdaContext.current().flatMap(ctx -> {
            ArdaEvent<T> event = new ArdaEvent<>(
                serviceName,
                eventType,
                data,
                ctx.traceId(),
                ctx.userId()
            );

            try {
                String payload = objectMapper.writeValueAsString(event);
                return Mono.fromFuture(kafkaTemplate.send(topic, event.id(), payload))
                    .then();
            } catch (JsonProcessingException e) {
                return Mono.error(e);
            }
        });
    }
}
