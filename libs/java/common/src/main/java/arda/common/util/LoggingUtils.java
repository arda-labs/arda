package arda.common.util;

import arda.common.context.ArdaContext;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.slf4j.MDC;
import reactor.core.publisher.Mono;
import reactor.core.publisher.Signal;

import java.util.function.Consumer;
import java.util.function.Function;

/**
 * Standardized Logging Utility for Arda Platform.
 * Supports Structured Logging and Reactive Context propagation.
 */
public class LoggingUtils {

    public static Logger getLogger(Class<?> clazz) {
        return LoggerFactory.getLogger(clazz);
    }

    /**
     * Log a message with context information from Reactor Context.
     * Usage: .flatMap(it -> LoggingUtils.logInfo(log, "Processing item: " + it).thenReturn(it))
     */
    public static Mono<Void> logInfo(Logger logger, String message) {
        return ArdaContext.current().map(ctx -> {
            MDC.put("userId", ctx.userId() != null ? ctx.userId() : "system");
            MDC.put("traceId", ctx.traceId() != null ? ctx.traceId() : "unknown");
            logger.info(message);
            MDC.clear();
            return (Void) null;
        });
    }

    /**
     * Helper to wrap log for Reactive pipelines using doOnEach.
     */
    public static <T> Consumer<Signal<T>> logOnNext(Logger logger, Function<T, String> logStatement) {
        return signal -> {
            if (signal.isOnNext()) {
                ArdaContext ctx = signal.getContextView().getOrDefault(ArdaContext.class, ArdaContext.empty());
                MDC.put("userId", ctx.userId() != null ? ctx.userId() : "system");
                MDC.put("traceId", ctx.traceId() != null ? ctx.traceId() : "unknown");
                logger.info(logStatement.apply(signal.get()));
                MDC.clear();
            }
        };
    }
}
