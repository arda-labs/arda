package arda.common.context;

import reactor.core.publisher.Mono;
import reactor.util.context.Context;
import java.util.Collections;
import java.util.List;

public record ArdaContext(
    String userId,
    String traceId,
    String tenantId,
    List<String> roles
) {
    private static final Class<ArdaContext> CONTEXT_KEY = ArdaContext.class;

    public ArdaContext {
        if (roles == null) roles = Collections.emptyList();
    }

    public static ArdaContext empty() {
        return new ArdaContext(null, null, null, Collections.emptyList());
    }

    public static Mono<ArdaContext> current() {
        return Mono.deferContextual(ctx -> {
            if (ctx.hasKey(CONTEXT_KEY)) {
                return Mono.just(ctx.get(CONTEXT_KEY));
            } else {
                return Mono.just(empty());
            }
        });
    }

    public static Context withContext(ArdaContext context) {
        return Context.of(CONTEXT_KEY, context);
    }

    public static Mono<String> getTraceId() {
        return current().map(ArdaContext::traceId);
    }
}
