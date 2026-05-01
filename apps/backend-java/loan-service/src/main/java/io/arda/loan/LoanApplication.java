package io.arda.loan;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication(scanBasePackages = {"io.arda.loan", "arda.messaging"})
public class LoanApplication {

    public static void main(String[] args) {
        // Đảm bảo JVM sử dụng múi giờ chuẩn để tránh lỗi "Asia/Saigon" trên Windows
        java.util.TimeZone.setDefault(java.util.TimeZone.getTimeZone("Asia/Ho_Chi_Minh"));
        System.setProperty("user.timezone", "Asia/Ho_Chi_Minh");

        SpringApplication.run(LoanApplication.class, args);
    }
}
