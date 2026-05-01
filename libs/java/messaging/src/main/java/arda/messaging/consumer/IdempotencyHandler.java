package arda.messaging.consumer;

import reactor.core.publisher.Mono;

/**
 * Interface for handling idempotency to prevent duplicate processing of messages.
 */
public interface IdempotencyHandler {
    /**
     * Check if the message has been processed.
     * @param key Unique key for the message (e.g. Event ID)
     * @return Mono<Boolean> true if already processed, false otherwise
     */
    Mono<Boolean> isProcessed(String key);

    /**
     * Mark the message as processed.
     */
    Mono<Void> markAsProcessed(String key);
}
