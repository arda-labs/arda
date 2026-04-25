package arda.messaging.event

import java.time.Instant
import java.util.UUID

/**
 * Standard Event structure for Arda Platform, based on CloudEvents spec.
 * This ensures consistency between Java and Go services when communicating via Redpanda/Kafka.
 */
data class ArdaEvent<T>(
    val id: String = UUID.randomUUID().toString(),
    val source: String, // Service name
    val type: String,   // Event type (e.g. arda.accounting.transaction.created)
    val data: T,
    val time: Instant = Instant.now(),
    val traceId: String? = null,
    val userId: String? = null
)
