package arda.messaging.consumer;

/**
 * Interface for handling idempotency to prevent duplicate processing of messages.
 * Synchronous and Virtual Thread safe.
 */
public interface IdempotencyHandler {
    /**
     * Check if the message has been processed.
     * @param key Unique key for the message (e.g. Event ID)
     * @return true if already processed, false otherwise
     */
    boolean isProcessed(String key);

    /**
     * Mark the message as processed.
     */
    void markAsProcessed(String key);
}
