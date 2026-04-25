package arda.accounting.application.service

import arda.accounting.application.dto.TransferRequest
import arda.accounting.domain.model.*
import arda.accounting.domain.service.AccountingDomainService
import arda.accounting.infrastructure.persistence.*
import arda.common.exception.BaseException
import arda.common.exception.ErrorCode
import arda.messaging.producer.ArdaEventProducer
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional
import reactor.core.publisher.Mono
import java.util.UUID

@Service
class AccountingApplicationService(
    private val accountRepository: AccountRepository,
    private val balanceRepository: BalanceRepository,
    private val journalRepository: JournalRepository,
    private val entryRepository: EntryRepository,
    private val eventProducer: ArdaEventProducer
) {
    private val domainService = AccountingDomainService()

    @Transactional
    fun transfer(request: TransferRequest): Mono<Journal> {
        return Mono.zip(
            accountRepository.findById(request.fromAccountId),
            accountRepository.findById(request.toAccountId),
            balanceRepository.findById(request.fromAccountId),
            balanceRepository.findById(request.toAccountId)
        ).flatMap { tuple ->
            val fromAccount = tuple.t1
            val toAccount = tuple.t2
            val fromBalance = tuple.t3
            val toBalance = tuple.t4

            // 1. Business Logic Validation
            if (fromBalance.availableBalance < request.amount) {
                return@flatMap Mono.error(BaseException(ErrorCode.INVALID_REQUEST, "Insufficient balance"))
            }

            // 2. Create Domain Journal
            val journal = domainService.createTransferJournal(
                fromAccount.toDomain(),
                toAccount.toDomain(),
                request.amount,
                request.description,
                request.referenceId
            )

            // 3. Persist Journal & Entries
            val journalEntity = journal.toEntity()
            val entryEntities = journal.entries.map { it.toEntity() }

            journalRepository.save(journalEntity)
                .thenMany(entryRepository.saveAll(entryEntities))
                // 4. Update Balances (Optimistic Locking handled by version check if needed, or simple update here)
                .then(balanceRepository.updateBalanceOptimistic(fromBalance.accountId, request.amount.negate(), fromBalance.version))
                .flatMap { updated ->
                    if (updated == 0) Mono.error(BaseException(ErrorCode.CONFLICT, "Concurrent update detected on sender account"))
                    else balanceRepository.updateBalanceOptimistic(toBalance.accountId, request.amount, toBalance.version)
                }
                .flatMap { updated ->
                    if (updated == 0) Mono.error(BaseException(ErrorCode.CONFLICT, "Concurrent update detected on receiver account"))
                    else Mono.just(journal)
                }
                // 5. Publish Event
                .flatMap { savedJournal ->
                    eventProducer.send(
                        topic = "accounting.journals",
                        eventType = "arda.accounting.journal.posted",
                        data = savedJournal,
                        serviceName = "accounting-service"
                    ).thenReturn(savedJournal)
                }
        }
    }

    // Mapper extensions
    private fun AccountEntity.toDomain() = Account(id, name, AccountType.valueOf(type), currency)
    private fun Journal.toEntity() = JournalEntity(id, description, referenceId, type.name, status.name)
    private fun Entry.toEntity() = EntryEntity(null, journalId, accountId, type.name, amount, currency, description)
}
