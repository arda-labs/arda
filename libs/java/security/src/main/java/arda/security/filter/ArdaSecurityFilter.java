package arda.security.filter;

import arda.common.context.ArdaContext;
import org.springframework.http.server.reactive.ServerHttpRequest;
import org.springframework.web.server.ServerWebExchange;
import org.springframework.web.server.WebFilter;
import org.springframework.web.server.WebFilterChain;
import reactor.core.publisher.Mono;

import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.stream.Collectors;

/**
 * Security Filter to extract user information from Gateway headers.
 * This follows the "Trust the Gateway" pattern used in large-scale systems.
 */
public class ArdaSecurityFilter implements WebFilter {

    public static final String HEADER_USER_ID = "X-User-Id";
    public static final String HEADER_TENANT_ID = "X-Tenant-Id";
    public static final String HEADER_ROLES = "X-User-Roles";
    public static final String HEADER_TRACE_ID = "X-Trace-Id";

    @Override
    public Mono<Void> filter(ServerWebExchange exchange, WebFilterChain chain) {
        ServerHttpRequest request = exchange.getRequest();

        ArdaContext context = new ArdaContext(
            request.getHeaders().getFirst(HEADER_USER_ID),
            request.getHeaders().getFirst(HEADER_TRACE_ID) != null ? request.getHeaders().getFirst(HEADER_TRACE_ID) : exchange.getRequest().getId(),
            request.getHeaders().getFirst(HEADER_TENANT_ID),
            extractRoles(request)
        );

        // Propagate the context to the Reactive stream
        return chain.filter(exchange)
            .contextWrite(ctx -> ArdaContext.withContext(context));
    }

    private List<String> extractRoles(ServerHttpRequest request) {
        String rolesHeader = request.getHeaders().getFirst(HEADER_ROLES);
        if (rolesHeader == null || rolesHeader.isBlank()) {
            return Collections.emptyList();
        }
        return Arrays.stream(rolesHeader.split(","))
            .map(String::trim)
            .collect(Collectors.toList());
    }
}
