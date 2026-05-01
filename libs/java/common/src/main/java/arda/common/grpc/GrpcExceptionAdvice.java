package arda.common.grpc;

import arda.common.exception.BaseException;
import io.grpc.Status;
import net.devh.boot.grpc.server.advice.GrpcAdvice;
import net.devh.boot.grpc.server.advice.GrpcExceptionHandler;

@GrpcAdvice
public class GrpcExceptionAdvice {

    @GrpcExceptionHandler(BaseException.class)
    public Status handleBaseException(BaseException e) {
        return Status.INTERNAL.withDescription(e.getMessage()).withCause(e);
    }

    @GrpcExceptionHandler(Exception.class)
    public Status handleException(Exception e) {
        return Status.UNKNOWN.withDescription("An unexpected error occurred").withCause(e);
    }
}
