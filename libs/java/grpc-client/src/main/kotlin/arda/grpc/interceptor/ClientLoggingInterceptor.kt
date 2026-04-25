package arda.grpc.interceptor

import io.grpc.*
import org.slf4j.LoggerFactory

class ClientLoggingInterceptor : ClientInterceptor {
    private val logger = LoggerFactory.getLogger(ClientLoggingInterceptor::class.java)

    override fun <ReqT, RespT> interceptCall(
        method: MethodDescriptor<ReqT, RespT>,
        callOptions: CallOptions,
        next: Channel
    ): ClientCall<ReqT, RespT> {
        logger.debug("gRPC Call: ${method.fullMethodName}")
        return next.newCall(method, callOptions)
    }
}
