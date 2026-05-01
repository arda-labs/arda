package arda.security.filter;

import arda.common.context.ArdaContext;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.web.filter.OncePerRequestFilter;

import java.io.IOException;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.stream.Collectors;

/**
 * Security Filter to extract user information from Gateway headers.
 * Synchronous and Virtual Thread safe using ArdaContext (ThreadLocal).
 */
public class ArdaSecurityFilter extends OncePerRequestFilter {

    public static final String HEADER_USER_ID = "X-User-Id";
    public static final String HEADER_TENANT_ID = "X-Tenant-Id";
    public static final String HEADER_ROLES = "X-User-Roles";
    public static final String HEADER_TRACE_ID = "X-Trace-Id";

    @Override
    protected void doFilterInternal(HttpServletRequest request, HttpServletResponse response, FilterChain filterChain)
            throws ServletException, IOException {

        ArdaContext context = new ArdaContext(
            request.getHeader(HEADER_USER_ID),
            request.getHeader(HEADER_TRACE_ID) != null ? request.getHeader(HEADER_TRACE_ID) : UUID_TRACE(),
            request.getHeader(HEADER_TENANT_ID),
            extractRoles(request)
        );

        try {
            // Set context for the current Virtual Thread
            ArdaContext.set(context);
            filterChain.doFilter(request, response);
        } finally {
            // CRITICAL: Clear context after request finishes to prevent leak
            ArdaContext.clear();
        }
    }

    private String UUID_TRACE() {
        return java.util.UUID.randomUUID().toString().replace("-", "");
    }

    private List<String> extractRoles(HttpServletRequest request) {
        String rolesHeader = request.getHeader(HEADER_ROLES);
        if (rolesHeader == null || rolesHeader.isBlank()) {
            return Collections.emptyList();
        }
        return Arrays.stream(rolesHeader.split(","))
            .map(String::trim)
            .collect(Collectors.toList());
    }
}
