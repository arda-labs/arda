package io.arda.crm.domain

import org.springframework.data.annotation.Id
import org.springframework.data.relational.core.mapping.Table
import java.util.UUID

@Table("customers")
data class Customer(
    @Id val id: UUID? = null,
    val customerCode: String,
    val name: String,
    val status: String, // PENDING, ACTIVE, REJECTED
    val cccdFileId: UUID?
)
