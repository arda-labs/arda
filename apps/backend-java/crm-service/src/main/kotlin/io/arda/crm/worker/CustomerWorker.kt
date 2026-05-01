package io.arda.crm.worker

import io.arda.crm.domain.CustomerRepository
import io.camunda.zeebe.client.api.response.ActivatedJob
import io.camunda.zeebe.spring.client.annotation.JobWorker
import kotlinx.coroutines.reactive.awaitFirstOrNull
import org.springframework.stereotype.Component
import java.util.UUID

@Component
class CustomerWorker(private val repository: CustomerRepository) {

    @JobWorker(type = "activate-customer")
    suspend fun activateCustomer(job: ActivatedJob) {
        val customerIdStr = job.variablesAsMap["customerId"] as? String ?: return
        val customerId = UUID.fromString(customerIdStr)

        val customer = repository.findById(customerId).awaitFirstOrNull()
        if (customer != null) {
            val updatedCustomer = customer.copy(status = "ACTIVE")
            repository.save(updatedCustomer).awaitFirstOrNull()
        }
    }
}
