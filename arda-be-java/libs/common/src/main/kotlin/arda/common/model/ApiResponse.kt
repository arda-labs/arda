package arda.common.model

import arda.common.exception.ErrorCode
import java.time.Instant

data class ApiResponse<T>(
    val success: Boolean,
    val data: T? = null,
    val error: ApiError? = null,
    val metadata: Map<String, Any>? = null,
    val timestamp: Instant = Instant.now()
) {
    companion object {
        fun <T> success(data: T?, metadata: Map<String, Any>? = null): ApiResponse<T> = ApiResponse(
            success = true,
            data = data,
            metadata = metadata
        )

        fun <T> error(errorCode: ErrorCode, message: String? = null, details: List<ErrorDetail>? = null): ApiResponse<T> = ApiResponse(
            success = false,
            error = ApiError(
                code = errorCode.code,
                message = message ?: errorCode.message,
                details = details
            )
        )
    }
}

data class ApiError(
    val code: String,
    val message: String,
    val details: List<ErrorDetail>? = null
)

data class ErrorDetail(
    val field: String? = null,
    val reason: String,
    val metadata: Map<String, Any>? = null
)
