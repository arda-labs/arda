package arda.grpc.interceptor;

import io.grpc.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class ClientLoggingInterceptor implements ClientInterceptor {
    private static final Logger logger = LoggerFactory.getLogger(ClientLoggingInterceptor.class);

    @Override
    public <ReqT, RespT> ClientCall<ReqT, RespT> interceptCall(
            MethodDescriptor<ReqT, RespT> method,
            CallOptions callOptions,
            Channel next) {
        logger.debug("gRPC Call: {}", method.getFullMethodName());
        return next.newCall(method, callOptions);
    }
}
