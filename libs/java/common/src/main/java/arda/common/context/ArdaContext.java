package arda.common.context;

import java.util.Collections;
import java.util.List;
import java.util.Optional;

public record ArdaContext(
    String userId,
    String traceId,
    String tenantId,
    List<String> roles
) {
    // ThreadLocal for Imperative/Virtual Threads flow
    private static final ThreadLocal<ArdaContext> THREAD_LOCAL = new ThreadLocal<>();

    public ArdaContext {
        if (roles == null) roles = Collections.emptyList();
    }

    public static ArdaContext empty() {
        return new ArdaContext(null, null, null, Collections.emptyList());
    }

    /**
     * For Imperative/Virtual Threads flow (Standard for Java 25+).
     */
    public static ArdaContext current() {
        return Optional.ofNullable(THREAD_LOCAL.get()).orElse(empty());
    }

    public static void set(ArdaContext context) {
        THREAD_LOCAL.set(context);
    }

    public static void clear() {
        THREAD_LOCAL.remove();
    }

    public static String getTraceId() {
        return Optional.ofNullable(current().traceId()).orElse("unknown");
    }
}
