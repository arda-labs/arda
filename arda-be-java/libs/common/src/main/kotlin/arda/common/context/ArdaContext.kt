package arda.common.context

import reactor.core.publisher.Mono
import reactor.util.context.Context

/**
 * Standardized Request Context for Arda Platform.
 * Follows Cloud-native patterns for context propagation in Reactive systems.
 */
data class ArdaContext(
    val userId: String? = null,
    val traceId: String? = null,
    val tenantId: String? = null,
    val roles: List<String> = emptyList()
) {
    companion object {
        private val CONTEXT_KEY = ArdaContext::class.java

        /**
         * Get the current ArdaContext from Reactor Context.
         */
        fun current(): Mono<ArdaContext> {
            return Mono.deferContextual { ctx ->
                if (ctx.hasKey(CONTEXT_KEY)) {
                    Mono.just(ctx.get(CONTEXT_KEY))
                } else {
                    Mono.just(ArdaContext())
                }
            }
        }

        /**
         * Create a Reactor Context with ArdaContext.
         */
        fun withContext(context: ArdaContext): Context {
            return Context.of(CONTEXT_KEY, context)
        }

        /**
         * Helper to get traceId for logging.
         */
        fun getTraceId(): Mono<String?> = current().map { it.traceId }
    }
}
