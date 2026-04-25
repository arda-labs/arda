package arda.accounting.domain.model

import java.math.BigDecimal
import java.time.Instant

data class Journal(
    val id: String,
    val description: String?,
    val referenceId: String?, // External ID (e.g., from EPAS)
    val type: JournalType,
    val status: JournalStatus = JournalStatus.POSTED,
    val entries: List<Entry> = emptyList(),
    val createdAt: Instant = Instant.now()
) {
    init {
        if (entries.isNotEmpty()) {
            validateDoubleEntry()
        }
    }

    private fun validateDoubleEntry() {
        val totalDebit = entries.filter { it.type == EntryType.DEBIT }.sumOf { it.amount }
        val totalCredit = entries.filter { it.type == EntryType.CREDIT }.sumOf { it.amount }
        require(totalDebit.compareTo(totalCredit) == 0) {
            "Invalid double-entry: Total Debit ($totalDebit) must equal Total Credit ($totalCredit)"
        }
    }
}

enum class JournalType {
    TRANSFER, DEPOSIT, WITHDRAW, FEE, ADJUSTMENT
}

enum class JournalStatus {
    POSTED, REVERSED, PENDING
}

data class Entry(
    val id: Long? = null,
    val journalId: String,
    val accountId: String,
    val type: EntryType,
    val amount: BigDecimal,
    val currency: String,
    val description: String?,
    val createdAt: Instant = Instant.now()
)

enum class EntryType {
    DEBIT, CREDIT
}
