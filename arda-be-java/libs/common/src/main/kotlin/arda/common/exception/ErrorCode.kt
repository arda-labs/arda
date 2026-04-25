package arda.common.exception

enum class ErrorCode(val code: String, val message: String) {
    INTERNAL_SERVER_ERROR("ERR_000", "Internal Server Error"),
    INVALID_REQUEST("ERR_001", "Invalid Request"),
    UNAUTHORIZED("ERR_002", "Unauthorized"),
    FORBIDDEN("ERR_003", "Forbidden"),
    NOT_FOUND("ERR_004", "Resource Not Found"),
    CONFLICT("ERR_005", "Resource Conflict"),
    EXTERNAL_SERVICE_ERROR("ERR_006", "External Service Error")
}
