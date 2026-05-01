package io.arda.hrm.domain;

import org.springframework.data.annotation.Id;
import org.springframework.data.relational.core.mapping.Table;
import java.util.UUID;

@Table("customers")
public record Customer(
    @Id UUID id,
    String customerCode,
    String name,
    String status, // PENDING, ACTIVE, REJECTED
    UUID cccdFileId
) {}
