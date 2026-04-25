package arda.accounting.application.dto

import java.math.BigDecimal

data class TransferRequest(
    val fromAccountId: String,
    val toAccountId: String,
    val amount: BigDecimal,
    val description: String?,
    val referenceId: String? = null
)
