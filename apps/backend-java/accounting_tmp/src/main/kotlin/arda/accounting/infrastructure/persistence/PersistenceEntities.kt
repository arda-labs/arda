package arda.accounting.infrastructure.persistence

import arda.database.entity.BaseEntity
import org.springframework.data.annotation.Id
import org.springframework.data.relational.core.mapping.Table
import java.math.BigDecimal
import java.time.Instant

@Table("accounts")
data class AccountEntity(
    @Id val id: String,
    val name: String,
    val type: String,
    val currency: String,
    val status: String
) : BaseEntity()

@Table("balances")
data class BalanceEntity(
    @Id val accountId: String,
    val availableBalance: BigDecimal,
    val lockedBalance: BigDecimal,
    val version: Long = 0,
    val updatedAt: Instant = Instant.now()
)

@Table("journals")
data class JournalEntity(
    @Id val id: String,
    val description: String?,
    val referenceId: String?,
    val type: String,
    val status: String,
    val createdAt: Instant = Instant.now(),
    val createdBy: String? = null
)

@Table("entries")
data class EntryEntity(
    @Id val id: Long? = null,
    val journalId: String,
    val accountId: String,
    val type: String,
    val amount: BigDecimal,
    val currency: String,
    val description: String?,
    val createdAt: Instant = Instant.now()
)
