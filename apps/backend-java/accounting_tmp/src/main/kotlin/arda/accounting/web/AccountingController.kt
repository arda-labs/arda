package arda.accounting.web

import arda.accounting.application.dto.TransferRequest
import arda.accounting.application.service.AccountingApplicationService
import arda.accounting.domain.model.Journal
import arda.common.model.ApiResponse
import arda.common.util.LoggingUtils
import org.springframework.web.bind.annotation.*
import reactor.core.publisher.Mono

@RestController
@RequestMapping("/api/v1/accounting")
class AccountingController(
    private val service: AccountingApplicationService
) {
    private val log = LoggingUtils.getLogger(AccountingController::class.java)

    @PostMapping("/transfer")
    fun transfer(@RequestBody request: TransferRequest): Mono<ApiResponse<Journal>> {
        return LoggingUtils.logInfo(log, "Received transfer request: $request")
            .then(service.transfer(request))
            .map { ApiResponse.success(it) }
            .doOnError { e -> log.error("Transfer failed: ${e.message}", e) }
    }
}
