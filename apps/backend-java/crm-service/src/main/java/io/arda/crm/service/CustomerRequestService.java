package io.arda.crm.service;

import arda.messaging.producer.ArdaEventProducer;
import io.arda.crm.domain.Customer;
import io.arda.crm.domain.CustomerRepository;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import java.util.Map;
import java.util.UUID;

@Service
public class CustomerRequestService {
    private final ArdaEventProducer eventProducer;
    private final CustomerRepository customerRepository;

    public CustomerRequestService(ArdaEventProducer eventProducer, CustomerRepository customerRepository) {
        this.eventProducer = eventProducer;
        this.customerRepository = customerRepository;
    }

    @Transactional
    public void createRegistrationRequest(Map<String, Object> customerData) {
        // 1. Lưu vào DB (Imperative style)
        UUID customerId = UUID.randomUUID();
        Customer customer = new Customer(
            customerId,
            (String) customerData.getOrDefault("customerCode", "CUST-" + customerId.toString().substring(0, 8)),
            (String) customerData.get("name"),
            "PENDING",
            customerData.containsKey("cccdFileId") ? UUID.fromString((String) customerData.get("cccdFileId")) : null
        );

        Customer saved = customerRepository.save(customer);

        // 2. Bắn event sang Kafka (Now synchronous)
        customerData.put("id", saved.getId().toString());
        eventProducer.send(
            "crm-events",
            "CUSTOMER_REGISTRATION_CREATED",
            customerData,
            "crm-service"
        );
    }
}
