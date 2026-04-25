package arda.security.filter

import arda.common.context.ArdaContext
import org.springframework.http.server.reactive.ServerHttpRequest
import org.springframework.web.server.ServerWebExchange
import org.springframework.web.server.WebFilter
import org.springframework.web.server.WebFilterChain
import reactor.core.publisher.Mono

/**
 * Security Filter to extract user information from Gateway headers.
 * This follows the "Trust the Gateway" pattern used in large-scale systems.
 */
class ArdaSecurityFilter : WebFilter {

    companion object {
        const val HEADER_USER_ID = "X-User-Id"
        const val HEADER_TENANT_ID = "X-Tenant-Id"
        const val HEADER_ROLES = "X-User-Roles"
        const val HEADER_TRACE_ID = "X-Trace-Id"
    }

    override fun filter(exchange: ServerWebExchange, chain: WebFilterChain): Mono<Void> {
        val request = exchange.request

        val context = ArdaContext(
            userId = request.headers.getFirst(HEADER_USER_ID),
            tenantId = request.headers.getFirst(HEADER_TENANT_ID),
            traceId = request.headers.getFirst(HEADER_TRACE_ID) ?: request.id,
            roles = extractRoles(request)
        )

        // Propagate the context to the Reactive stream
        return chain.filter(exchange)
            .contextWrite { ArdaContext.withContext(context) }
    }

    private fun extractRoles(request: ServerHttpRequest): List<String> {
        val rolesHeader = request.headers.getFirst(HEADER_ROLES)
        return rolesHeader?.split(",")?.map { it.trim() } ?: emptyList()
    }
}
