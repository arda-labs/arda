package arda.common.exception;

import java.util.Map;

public class BaseException extends RuntimeException {
    private final ErrorCode errorCode;
    private final Map<String, Object> details;

    public BaseException(ErrorCode errorCode) {
        this(errorCode, errorCode.getMessage(), null, null);
    }

    public BaseException(ErrorCode errorCode, String message) {
        this(errorCode, message, null, null);
    }

    public BaseException(ErrorCode errorCode, String message, Map<String, Object> details) {
        this(errorCode, message, details, null);
    }

    public BaseException(ErrorCode errorCode, String message, Map<String, Object> details, Throwable cause) {
        super(message, cause);
        this.errorCode = errorCode;
        this.details = details;
    }

    public ErrorCode getErrorCode() {
        return errorCode;
    }

    public Map<String, Object> getDetails() {
        return details;
    }
}
