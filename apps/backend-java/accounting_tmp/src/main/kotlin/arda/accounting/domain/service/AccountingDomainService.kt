package arda.accounting.domain.service

import arda.accounting.domain.model.*
import arda.common.exception.BaseException
import arda.common.exception.ErrorCode
import java.math.BigDecimal
import java.util.UUID

class AccountingDomainService {

    /**
     * Creates a double-entry journal for a transfer between two accounts.
     */
    fun createTransferJournal(
        fromAccount: Account,
        toAccount: Account,
        amount: BigDecimal,
        description: String?,
        referenceId: String?
    ): Journal {
        val journalId = UUID.randomUUID().toString()

        val debitEntry = Entry(
            journalId = journalId,
            accountId = fromAccount.id,
            type = EntryType.DEBIT,
            amount = amount,
            currency = fromAccount.currency,
            description = "Debit for transfer to ${toAccount.id}"
        )

        val creditEntry = Entry(
            journalId = journalId,
            accountId = toAccount.id,
            type = EntryType.CREDIT,
            amount = amount,
            currency = toAccount.currency,
            description = "Credit for transfer from ${fromAccount.id}"
        )

        return Journal(
            id = journalId,
            description = description,
            referenceId = referenceId,
            type = JournalType.TRANSFER,
            entries = listOf(debitEntry, creditEntry)
        )
    }
}
