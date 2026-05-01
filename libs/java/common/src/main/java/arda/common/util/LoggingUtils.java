package arda.common.util;

import arda.common.context.ArdaContext;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.slf4j.MDC;

import java.util.function.Supplier;

/**
 * Standardized Logging Utility for Arda Platform.
 * Supports Structured Logging and Virtual Thread context propagation.
 */
public class LoggingUtils {

    public static Logger getLogger(Class<?> clazz) {
        return LoggerFactory.getLogger(clazz);
    }

    /**
     * Log a message with context information.
     * Synchronous and Virtual Thread safe.
     */
    public static void info(Logger logger, String message) {
        ArdaContext ctx = ArdaContext.current();
        try {
            MDC.put("userId", ctx.userId() != null ? ctx.userId() : "system");
            MDC.put("traceId", ctx.traceId() != null ? ctx.traceId() : "unknown");
            MDC.put("tenantId", ctx.tenantId() != null ? ctx.tenantId() : "shared");
            logger.info(message);
        } finally {
            MDC.clear();
        }
    }

    /**
     * Executes a supplier within a logging context.
     */
    public static <T> T withContext(Supplier<T> supplier) {
        ArdaContext ctx = ArdaContext.current();
        try {
            MDC.put("userId", ctx.userId() != null ? ctx.userId() : "system");
            MDC.put("traceId", ctx.traceId() != null ? ctx.traceId() : "unknown");
            return supplier.get();
        } finally {
            MDC.clear();
        }
    }
}
