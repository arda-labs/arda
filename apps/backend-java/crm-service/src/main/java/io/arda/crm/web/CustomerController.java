package io.arda.crm.web;

import io.arda.crm.service.CustomerRequestService;
import org.springframework.web.bind.annotation.*;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/customers")
public class CustomerController {
    private final CustomerRequestService requestService;

    public CustomerController(CustomerRequestService requestService) {
        this.requestService = requestService;
    }

    @PostMapping("/register")
    public Map<String, Object> register(@RequestBody Map<String, Object> data) {
        requestService.createRegistrationRequest(data);
        return Map.of(
            "success", true,
            "message", "Registration request submitted successfully",
            "id", data.getOrDefault("id", "UNKNOWN")
        );
    }
}
