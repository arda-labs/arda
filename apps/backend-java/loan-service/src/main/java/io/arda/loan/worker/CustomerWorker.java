package io.arda.loan.worker;

import io.arda.loan.domain.Customer;
import io.arda.loan.domain.CustomerRepository;
import io.camunda.zeebe.client.api.response.ActivatedJob;
import io.camunda.zeebe.spring.client.annotation.JobWorker;
import org.springframework.stereotype.Component;
import reactor.core.publisher.Mono;

import java.util.UUID;

@Component
public class CustomerWorker {
    private final CustomerRepository repository;

    public CustomerWorker(CustomerRepository repository) {
        this.repository = repository;
    }

    @JobWorker(type = "activate-customer")
    public Mono<Void> activateCustomer(ActivatedJob job) {
        String customerIdStr = (String) job.getVariablesAsMap().get("customerId");
        if (customerIdStr == null) {
            return Mono.empty();
        }

        UUID customerId = UUID.fromString(customerIdStr);

        return repository.findById(customerId)
                .flatMap(customer -> {
                    Customer updatedCustomer = new Customer(
                        customer.id(),
                        customer.customerCode(),
                        customer.name(),
                        "ACTIVE",
                        customer.cccdFileId()
                    );
                    return repository.save(updatedCustomer);
                })
                .then();
    }
}
