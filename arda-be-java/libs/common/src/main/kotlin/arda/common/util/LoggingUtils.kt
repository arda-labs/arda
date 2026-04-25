package arda.common.util

import arda.common.context.ArdaContext
import org.slf4j.Logger
import org.slf4j.LoggerFactory
import org.slf4j.MDC
import reactor.core.publisher.Mono
import reactor.core.publisher.Signal

/**
 * Standardized Logging Utility for Arda Platform.
 * Supports Structured Logging and Reactive Context propagation.
 */
object LoggingUtils {

    fun getLogger(clazz: Class<*>): Logger = LoggerFactory.getLogger(clazz)

    /**
     * Log a message with context information from Reactor Context.
     * Usage: .doOnNext { LoggingUtils.logInfo(log, "Processing item: $it") }
     */
    fun logInfo(logger: Logger, message: String): Mono<Unit> {
        return ArdaContext.current().map { ctx ->
            MDC.put("userId", ctx.userId ?: "system")
            MDC.put("traceId", ctx.traceId ?: "unknown")
            logger.info(message)
            MDC.clear()
        }
    }

    /**
     * Helper to wrap log for Reactive pipelines using doOnEach.
     */
    fun <T> logOnNext(logger: Logger, logStatement: (T) -> String): (Signal<T>) -> Unit {
        return { signal ->
            if (signal.isOnNext) {
                val ctx = signal.contextView.getOrDefault(ArdaContext::class.java, ArdaContext())
                MDC.put("userId", ctx.userId ?: "system")
                MDC.put("traceId", ctx.traceId ?: "unknown")
                logger.info(logStatement(signal.get()!!))
                MDC.clear()
            }
        }
    }
}
