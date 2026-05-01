package arda.grpc.interceptor;

import arda.common.context.ArdaContext;
import io.grpc.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Interceptor to propagate ArdaContext (traceId, userId) to downstream gRPC services.
 * Ensures consistency between Java and Go services.
 */
public class ClientContextInterceptor implements ClientInterceptor {
    private static final Logger logger = LoggerFactory.getLogger(ClientContextInterceptor.class);

    public static final Metadata.Key<String> USER_ID_KEY = Metadata.Key.of("x-user-id", Metadata.ASCII_STRING_MARSHALLER);
    public static final Metadata.Key<String> TRACE_ID_KEY = Metadata.Key.of("x-trace-id", Metadata.ASCII_STRING_MARSHALLER);

    @Override
    public <ReqT, RespT> ClientCall<ReqT, RespT> interceptCall(
            MethodDescriptor<ReqT, RespT> method,
            CallOptions callOptions,
            Channel next) {
        return new ForwardingClientCall.SimpleForwardingClientCall<>(
                next.newCall(method, callOptions)
        ) {
            @Override
            public void start(Listener<RespT> responseListener, Metadata headers) {
                // Get context from Reactor and propagate to Metadata
                ArdaContext.current().subscribe(ctx -> {
                    if (ctx.userId() != null) headers.put(USER_ID_KEY, ctx.userId());
                    if (ctx.traceId() != null) headers.put(TRACE_ID_KEY, ctx.traceId());
                });
                super.start(responseListener, headers);
            }
        };
    }
}
