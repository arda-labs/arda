package io.arda.crm.worker;

import io.arda.crm.domain.Customer;
import io.arda.crm.domain.CustomerRepository;
import io.camunda.zeebe.client.api.response.ActivatedJob;
import io.camunda.zeebe.spring.client.annotation.JobWorker;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

@Component
public class CustomerWorker {
    private static final Logger log = LoggerFactory.getLogger(CustomerWorker.class);
    private final CustomerRepository repository;

    public CustomerWorker(CustomerRepository repository) {
        this.repository = repository;
    }

    @JobWorker(type = "activate-customer")
    @Transactional
    public void activateCustomer(final ActivatedJob job) {
        var variables = job.getVariablesAsMap();
        var customerIdStr = (String) variables.get("customerId");

        if (customerIdStr == null) {
            log.warn("Job {}: Missing customerId variable", job.getKey());
            return;
        }

        log.info("Zeebe Worker: Activating customer {}", customerIdStr);
        UUID customerId = UUID.fromString(customerIdStr);

        repository.findById(customerId).ifPresentOrElse(customer -> {
            customer.setStatus("ACTIVE");
            repository.save(customer);
            log.info("Customer {} activated successfully", customerId);
        }, () -> {
            log.error("Customer {} not found for activation", customerId);
        });
    }
}
