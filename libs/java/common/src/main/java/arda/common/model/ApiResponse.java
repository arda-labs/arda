package arda.common.model;

import arda.common.exception.ErrorCode;
import java.time.Instant;
import java.util.List;
import java.util.Map;

public record ApiResponse<T>(
    boolean success,
    T data,
    ApiError error,
    Map<String, Object> metadata,
    Instant timestamp
) {
    public ApiResponse {
        if (timestamp == null) timestamp = Instant.now();
    }

    public static <T> ApiResponse<T> success(T data) {
        return new ApiResponse<>(true, data, null, null, Instant.now());
    }

    public static <T> ApiResponse<T> success(T data, Map<String, Object> metadata) {
        return new ApiResponse<>(true, data, null, metadata, Instant.now());
    }

    public static <T> ApiResponse<T> error(ErrorCode errorCode) {
        return new ApiResponse<>(false, null, new ApiError(errorCode.getCode(), errorCode.getMessage(), null), null, Instant.now());
    }

    public static <T> ApiResponse<T> error(ErrorCode errorCode, String message) {
        return new ApiResponse<>(false, null, new ApiError(errorCode.getCode(), message, null), null, Instant.now());
    }

    public record ApiError(String code, String message, List<ErrorDetail> details) {}

    public record ErrorDetail(String field, String reason, Map<String, Object> metadata) {}
}
