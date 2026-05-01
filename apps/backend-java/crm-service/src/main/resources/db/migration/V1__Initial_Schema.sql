-- V1__Initial_Schema.sql
CREATE TABLE customers (
    id UUID PRIMARY KEY,
    customer_code VARCHAR(255),
    name VARCHAR(255),
    status VARCHAR(50), -- PENDING, ACTIVE, REJECTED
    cccd_file_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    updated_by VARCHAR(255)
);
