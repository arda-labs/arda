package io.arda.crm.domain;

import arda.database.entity.BaseEntity;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import java.util.UUID;

@Entity
@Table(name = "customers")
public class Customer extends BaseEntity {
    @Id
    private UUID id;
    private String customerCode;
    private String name;
    private String status; // PENDING, ACTIVE, REJECTED
    private UUID cccdFileId;

    public Customer() {}

    public Customer(UUID id, String customerCode, String name, String status, UUID cccdFileId) {
        this.id = id;
        this.customerCode = customerCode;
        this.name = name;
        this.status = status;
        this.cccdFileId = cccdFileId;
    }

    // Getters and Setters
    public UUID getId() { return id; }
    public void setId(UUID id) { this.id = id; }
    public String getCustomerCode() { return customerCode; }
    public void setCustomerCode(String customerCode) { this.customerCode = customerCode; }
    public String getName() { return name; }
    public void setName(String name) { this.name = name; }
    public String getStatus() { return status; }
    public void setStatus(String status) { this.status = status; }
    public UUID getCccdFileId() { return cccdFileId; }
    public void setCccdFileId(UUID cccdFileId) { this.cccdFileId = cccdFileId; }
}
