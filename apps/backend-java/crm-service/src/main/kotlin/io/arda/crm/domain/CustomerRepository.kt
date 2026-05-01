package io.arda.crm.domain

import org.springframework.data.repository.reactive.ReactiveCrudRepository
import org.springframework.stereotype.Repository
import java.util.UUID

@Repository
interface CustomerRepository : ReactiveCrudRepository<Customer, UUID>
