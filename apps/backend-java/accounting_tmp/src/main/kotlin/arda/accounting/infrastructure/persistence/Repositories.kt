package arda.accounting.infrastructure.persistence

import org.springframework.data.r2dbc.repository.Modifying
import org.springframework.data.r2dbc.repository.Query
import org.springframework.data.repository.reactive.ReactiveCrudRepository
import reactor.core.publisher.Flux
import reactor.core.publisher.Mono
import java.math.BigDecimal

interface AccountRepository : ReactiveCrudRepository<AccountEntity, String>

interface BalanceRepository : ReactiveCrudRepository<BalanceEntity, String> {

    @Modifying
    @Query("UPDATE balances SET available_balance = available_balance + :amount, updated_at = NOW(), version = version + 1 WHERE account_id = :accountId AND version = :version")
    fun updateBalanceOptimistic(accountId: String, amount: BigDecimal, version: Long): Mono<Int>
}

interface JournalRepository : ReactiveCrudRepository<JournalEntity, String>

interface EntryRepository : ReactiveCrudRepository<EntryEntity, Long> {
    fun findByJournalId(journalId: String): Flux<EntryEntity>
    fun findByAccountId(accountId: String): Flux<EntryEntity>
}
