package io.arda.crm.grpc;

import io.arda.crm.domain.CustomerRepository;
import io.arda.crm.v1.CRMGrpc;
import io.arda.crm.v1.FinalizeCustomerReply;
import io.arda.crm.v1.FinalizeCustomerRequest;
import io.grpc.stub.StreamObserver;
import net.devh.boot.grpc.server.service.GrpcService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.transaction.annotation.Transactional;
import java.util.UUID;

@GrpcService
public class CrmGrpcService extends CRMGrpc.CRMImplBase {
    private static final Logger log = LoggerFactory.getLogger(CrmGrpcService.class);
    private final CustomerRepository customerRepository;

    public CrmGrpcService(CustomerRepository customerRepository) {
        this.customerRepository = customerRepository;
    }

    @Override
    @Transactional
    public void finalizeCustomer(FinalizeCustomerRequest request, StreamObserver<FinalizeCustomerReply> responseObserver) {
        log.info("gRPC Server (Imperative): Finalizing customer {} with status {}", request.getCustomerId(), request.getStatus());

        try {
            UUID customerId = UUID.fromString(request.getCustomerId());
            String status = request.getStatus();

            customerRepository.findById(customerId).ifPresentOrElse(customer -> {
                customer.setStatus(status);
                customerRepository.save(customer);

                FinalizeCustomerReply reply = FinalizeCustomerReply.newBuilder()
                    .setSuccess(true)
                    .setMessage("Customer " + request.getCustomerId() + " has been finalized with status: " + status)
                    .build();

                responseObserver.onNext(reply);
                responseObserver.onCompleted();
            }, () -> {
                responseObserver.onError(io.grpc.Status.NOT_FOUND
                    .withDescription("Customer not found: " + request.getCustomerId())
                    .asRuntimeException());
            });
        } catch (Exception e) {
            log.error("Error finalizing customer", e);
            responseObserver.onError(io.grpc.Status.INTERNAL
                .withDescription("Failed to finalize customer: " + e.getMessage())
                .asRuntimeException());
        }
    }
}
