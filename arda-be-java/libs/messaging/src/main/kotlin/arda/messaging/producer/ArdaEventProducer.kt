package arda.messaging.producer

import arda.common.context.ArdaContext
import arda.messaging.event.ArdaEvent
import com.fasterxml.jackson.databind.ObjectMapper
import org.springframework.kafka.core.KafkaTemplate
import org.springframework.stereotype.Component
import reactor.core.publisher.Mono

@Component
class ArdaEventProducer(
    private val kafkaTemplate: KafkaTemplate<String, String>,
    private val objectMapper: ObjectMapper
) {
    /**
     * Send an event with automatic Context propagation (traceId, userId).
     */
    fun <T> send(topic: String, eventType: String, data: T, serviceName: String): Mono<Void> {
        return ArdaContext.current().flatMap { ctx ->
            val event = ArdaEvent(
                source = serviceName,
                type = eventType,
                data = data,
                traceId = ctx.traceId,
                userId = ctx.userId
            )

            val payload = objectMapper.writeValueAsString(event)

            Mono.fromFuture(kafkaTemplate.send(topic, event.id, payload))
                .then()
        }
    }
}
