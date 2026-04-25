package arda.messaging.consumer

import reactor.core.publisher.Mono

/**
 * Interface for handling idempotency to prevent duplicate processing of messages.
 */
interface IdempotencyHandler {
    /**
     * Check if the message has been processed.
     * @param key Unique key for the message (e.g. Event ID)
     * @return Mono<Boolean> true if already processed, false otherwise
     */
    fun isProcessed(key: String): Mono<Boolean>

    /**
     * Mark the message as processed.
     */
    fun markAsProcessed(key: String): Mono<Void>
}
