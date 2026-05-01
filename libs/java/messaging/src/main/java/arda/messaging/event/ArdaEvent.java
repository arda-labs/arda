package arda.messaging.event;

import java.time.Instant;
import java.util.UUID;

/**
 * Standard Event structure for Arda Platform, based on CloudEvents spec.
 * This ensures consistency between Java and Go services when communicating via Redpanda/Kafka.
 */
public record ArdaEvent<T>(
    String id,
    String source, // Service name
    String type,   // Event type (e.g. arda.accounting.transaction.created)
    T data,
    Instant time,
    String traceId,
    String userId
) {
    public ArdaEvent(String source, String type, T data, String traceId, String userId) {
        this(UUID.randomUUID().toString(), source, type, data, Instant.now(), traceId, userId);
    }
}
