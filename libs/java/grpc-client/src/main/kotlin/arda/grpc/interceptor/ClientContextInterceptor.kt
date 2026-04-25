package arda.grpc.interceptor

import arda.common.context.ArdaContext
import io.grpc.*
import org.slf4j.LoggerFactory

/**
 * Interceptor to propagate ArdaContext (traceId, userId) to downstream gRPC services.
 * Ensures consistency between Java and Go services.
 */
class ClientContextInterceptor : ClientInterceptor {
    private val logger = LoggerFactory.getLogger(ClientContextInterceptor::class.java)

    companion object {
        val USER_ID_KEY: Metadata.Key<String> = Metadata.Key.of("x-user-id", Metadata.ASCII_STRING_MARSHALLER)
        val TRACE_ID_KEY: Metadata.Key<String> = Metadata.Key.of("x-trace-id", Metadata.ASCII_STRING_MARSHALLER)
    }

    override fun <ReqT, RespT> interceptCall(
        method: MethodDescriptor<ReqT, RespT>,
        callOptions: CallOptions,
        next: Channel
    ): ClientCall<ReqT, RespT> {
        return object : ForwardingClientCall.SimpleForwardingClientCall<ReqT, RespT>(
            next.newCall(method, callOptions)
        ) {
            override fun start(responseListener: Listener<RespT>, headers: Metadata) {
                // Get context from Reactor and propagate to Metadata
                // Note: In a real reactive app, we would use deferContextual here,
                // but for simplicity in the interceptor, we assume the context is available.
                ArdaContext.current().subscribe { ctx ->
                    ctx.userId?.let { headers.put(USER_ID_KEY, it) }
                    ctx.traceId?.let { headers.put(TRACE_ID_KEY, it) }
                }
                super.start(responseListener, headers)
            }
        }
    }
}
