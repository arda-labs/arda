package arda.accounting.domain.model

import java.math.BigDecimal
import java.time.Instant

data class Account(
    val id: String,
    val name: String,
    val type: AccountType,
    val currency: String = "VND",
    val status: AccountStatus = AccountStatus.ACTIVE,
    val createdAt: Instant = Instant.now()
)

enum class AccountType {
    ASSET, LIABILITY, EQUITY, REVENUE, EXPENSE
}

enum class AccountStatus {
    ACTIVE, LOCKED, CLOSED
}

data class Balance(
    val accountId: String,
    var availableBalance: BigDecimal,
    var lockedBalance: BigDecimal,
    val version: Long = 0,
    val updatedAt: Instant = Instant.now()
) {
    fun hasEnoughBalance(amount: BigDecimal): Boolean {
        return availableBalance >= amount
    }

    fun credit(amount: BigDecimal): Balance {
        return copy(
            availableBalance = availableBalance.add(amount),
            updatedAt = Instant.now()
        )
    }

    fun debit(amount: BigDecimal): Balance {
        require(hasEnoughBalance(amount)) { "Insufficient balance" }
        return copy(
            availableBalance = availableBalance.subtract(amount),
            updatedAt = Instant.now()
        )
    }
}
