package arda.common.exception

open class BaseException(
    val errorCode: ErrorCode,
    override val message: String = errorCode.message,
    val details: Map<String, Any>? = null,
    cause: Throwable? = null
) : RuntimeException(message, cause)
