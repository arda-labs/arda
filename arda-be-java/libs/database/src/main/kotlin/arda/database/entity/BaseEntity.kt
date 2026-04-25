package arda.database.entity

import org.springframework.data.annotation.CreatedBy
import org.springframework.data.annotation.CreatedDate
import org.springframework.data.annotation.LastModifiedBy
import org.springframework.data.annotation.LastModifiedDate
import org.springframework.data.relational.core.mapping.Column
import java.time.Instant

abstract class BaseEntity {
    @CreatedDate
    @Column("created_at")
    var createdAt: Instant? = null

    @LastModifiedDate
    @Column("updated_at")
    var updatedAt: Instant? = null

    @CreatedBy
    @Column("created_by")
    var createdBy: String? = null

    @LastModifiedBy
    @Column("updated_by")
    var updatedBy: String? = null
}
